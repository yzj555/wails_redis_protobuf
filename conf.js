window.CONF_DATA = {
    "redisServer": [
        {
            "name": "本地服",
            "cluster": true,
            "address": "127.0.0.1:6379",
            "password": ""
        },
        {
            "name": "足小内网测试服",
            "cluster": false,
            "address": "192.168.8.158:6379",
            "password": ""
        },
        {
            "name": "足小QA服",
            "cluster": true,
            "address": "192.168.8.166:7000",
            "password": ""
        }
    ],
    "proto": {
        "dir": "E:\\football\\server\\server_v2\\server\\captain_tsubasa_server\\src\\idl\\pb"
    }
}