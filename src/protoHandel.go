package server

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/jhump/protoreflect/desc/protoprint"
	"github.com/jhump/protoreflect/dynamic"
	"github.com/vmihailenco/msgpack/v5"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var protoDescList = make([]*desc.FileDescriptor, 0)
var msgNameList = make([]string, 0)

// LoadProtoFiles 从 .proto 文件加载描述符
func LoadProtoFiles(dir string) {
	var filenames []string
	pathList := make([]string, 0)
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err // 处理遍历错误
		}
		if info.IsDir() {
			if strings.Contains(path, "google") || strings.Contains(path, ".idea") {
				return filepath.SkipDir
			}
			pathList = append(pathList, path)
			//fmt.Println("Directory:", path) // 输出目录路径
		} else {
			if filepath.Ext(info.Name()) == ".proto" {
				filenames = append(filenames, info.Name())
				//fmt.Println("File:", path) // 输出文件路径
			}
		}
		return nil
	})
	fmt.Println(filenames)
	parser := protoparse.Parser{ImportPaths: pathList}
	for _, filename := range filenames {
		desc, err := parser.ParseFiles(filename)
		if err != nil {
			fmt.Printf("ParseFiles err, %v\n", err)
			continue
		}
		protoDescList = append(protoDescList, desc...)
	}
	for _, desc := range protoDescList {
		for _, msg := range desc.GetMessageTypes() {
			//fmt.Println(msg.GetName())
			msgNameList = append(msgNameList, msg.GetName())
		}
	}
}

func UnmarshalByMsgpack(data []byte) interface{} {
	var res interface{}
	err := msgpack.Unmarshal(data, &res)
	if err != nil {
		fmt.Printf("UnmarshalByMsgpack err, %v\n", err)
		return data
	}
	return res
}

func UnmarshalByProto(data []byte, msgName string) interface{} {
	resStr := AutoUnmarshal([]byte(data), msgName)
	if resStr != "" {
		res := make(map[string]interface{})
		err := json.Unmarshal([]byte(resStr), &res)
		if err != nil {
			fmt.Printf("UnmarshalByProto err, %v\n", err)
			return nil
		}
		return res
	}
	fmt.Printf("AutoUnmarshal err 解析失败，尝试手动解析, %v\n", msgName)
	fields, err := ParseProtobufToJSON([]byte(data))
	if err != nil {
		fmt.Printf("ParseProtobufToJSON err, %v\n", err)
		return data
	} else {
		return fields
	}
}

func AutoUnmarshal(data []byte, msgName string) string {
	if msgName != "" {
		msgName = fmt.Sprintf("pb.%s", msgName)
	}
	for _, desc := range protoDescList {
		printer := &protoprint.Printer{}
		var buf bytes.Buffer
		printer.PrintProtoFile(desc, &buf)
		//err := printer.PrintProtoFile(desc, &buf)
		//if err != nil {
		//	fmt.Printf("PrintProtoFile err, %v\n", err)
		//	return ""
		//}
		//fmt.Printf("FileDescriptor: %s\n", buf.String())
		//通过proto的message名称得到MessageDescriptor 结构体定义描述符
		msg := desc.FindMessage(msgName)
		if msg != nil {
			//再用消息描述符，动态的构造一个pb消息体
			dMsg := dynamic.NewMessage(msg)
			dMsg.Unmarshal(data)
			//把 消息体序列化成 JSON 数据
			jsStr, _ := dMsg.MarshalJSONIndent()
			fmt.Printf("%s: %s\n", msgName, jsStr)
			return string(jsStr)
		}
	}
	return ""
}

func ParseProtobufToJSON(data []byte) (map[int]interface{}, error) {
	buf := bytes.NewReader(data)
	fields := make(map[int]interface{})

	for buf.Len() > 0 {
		// 读取 Tag (field number + wire type)
		tag, err := binary.ReadUvarint(buf)
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, fmt.Errorf("读取 Tag 失败: %v", err)
		}

		fieldNumber := int(tag >> 3)
		wireType := int(tag & 0x07)

		// 根据 Wire Type 解析值
		var value interface{}
		switch wireType {
		case 0: // Varint (int32, int64, bool, enum)
			val, err := binary.ReadUvarint(buf)
			if err != nil {
				return nil, fmt.Errorf("解析 Varint 失败: %v", err)
			}
			value = val

		case 1: // 64-bit (fixed64, double)
			var val uint64
			if err := binary.Read(buf, binary.LittleEndian, &val); err != nil {
				return nil, fmt.Errorf("解析 64-bit 失败: %v", err)
			}
			value = val

		case 2: // Length-delimited (string, bytes, 嵌套消息, 数组)
			length, err := binary.ReadUvarint(buf)
			if err != nil {
				return nil, fmt.Errorf("读取长度失败: %v", err)
			}
			bytesData := make([]byte, length)
			if _, err := io.ReadFull(buf, bytesData); err != nil {
				return nil, fmt.Errorf("读取字节失败: %v", err)
			}

			// 尝试递归解析嵌套消息
			if nestedFields, err := ParseProtobufToJSON(bytesData); err == nil {
				value = nestedFields
			} else {
				value = bytesData // 作为原始字节数组
			}

		case 5: // 32-bit (fixed32, float)
			var val uint32
			if err := binary.Read(buf, binary.LittleEndian, &val); err != nil {
				return nil, fmt.Errorf("解析 32-bit 失败: %v", err)
			}
			value = val

		default:
			return nil, fmt.Errorf("不支持的 Wire Type: %d", wireType)
		}
		fields[fieldNumber] = value
	}
	return fields, nil
}

func GetAllMsgName() []string {
	return msgNameList
}
