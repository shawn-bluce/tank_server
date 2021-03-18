package utils

import (
	"TankServer/zinx/ziface"
	"encoding/json"
	"io/ioutil"
)

/*
	存储一切有关zinx框架的全局参数，供其他模块使用
	一些参数是可以通过zinx.json由用户进行配置
*/
type GlobalObj struct {
	/*
		Server
	 */
	TcpServer ziface.IServer // 当前Zinx全局的Server对象
	Host string				 // 当前服务器主机监听的ip
	TcpPort int				 // 当前服务器主机监听的端口号
	Name string				 // 当前服务器的名称

	/*
		Zinx
	 */
	Version string			 // 当前Zinx的版本号
	MaxConn int				 // 当前服务器主机允许的最大连接数
	MaxPackageSize uint32	 // 当前Zinx框架数据包的最大值
	WorkerPoolSize uint32	 // 当前业务工作Worker池的Goroutine
	MaxWorkerTaskLen uint32  // zinx框架允许用户最多开辟多少个Worker
}

/*
	定义一个全局的对外GlobalObj
 */
var GlobalObject *GlobalObj


/*
	加载配置文件方法
 */
func (g *GlobalObj) Reload() {
	data, err := ioutil.ReadFile("F:\\Projects\\zinxDemo\\src\\mmo_game_zinx\\conf\\zinx.json")
	if err != nil {
		panic(err)
	}
	// 将json文件数据解析到struct中
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}


/*
	提供一个init方法，初始化当前的GlobalObject
	以下为默认配置文件，若用户有额外导入配置文件，则以用户导入为主
 */
func init()  {
	GlobalObject = &GlobalObj{
		Name: "ZinxServerApp",
		Version: "V0.8",
		TcpPort: 8999,
		Host: "0.0.0.0",
		MaxConn: 1000,
		MaxPackageSize: 4096,
		MaxWorkerTaskLen: 1024,
		WorkerPoolSize: 12,
	}

	// 应该尝试从conf/zinx.json中加载用户自定义的参数
	GlobalObject.Reload()
}
