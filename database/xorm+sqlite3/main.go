package main

import (
	"fmt"
	_ "os"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"xorm.io/xorm"
)

var engine *xorm.Engine

type User struct {
	Id        int64     `xorm: autoincr`
	Name      string    `xorm:"varchar(25) notnull 'usr_name'"`
	CreatedAt time.Time `xorm:"created"`
}

// https://xorm.io/zh/docs/chapter-09/readme/
func main() {

	// var myMap = make(map[string]string)
	// myMap["1"] = "aaaa"
	// fmt.Printf("xxxx%v\n", myMap)
	var myArray = make([]int, 2)
	myArray = append(myArray, 1)
	fmt.Println(cap(myArray))
	myArray = append(myArray, 1)
	myArray = append(myArray, 1)

	fmt.Println(cap(myArray))

	myArray = append(myArray, 1)
	myArray = append(myArray, 1)
	myArray = append(myArray, 1)
	myArray = append(myArray, 1)

	fmt.Println(cap(myArray))

	name := "GO"
	fmt.Println("xxx", name)

	var err error
	engine, err = xorm.NewEngine("sqlite3", "./test.db")
	// TODO: 貌似无效
	engine.TZLocation, _ = time.LoadLocation("Asia/Shanghai")
	fmt.Println(err)
	err = engine.Sync2(new(User))
	fmt.Println(err)

	// 删除指定数据
	engine.ID(1).Delete(new(User))

	// 删除整张表数据
	// 注意：像下面这样，删除整个表的数据，如果user中包含有bool,float64或者float32类型，有可能会报一个保护性的错误
	// engine.Delete(new(User))
	// 这个时候必须这样
	engine.Where("1=1").Delete(new(User))

	// 插入一条数据
	user := new(User)
	user.Name = "myname"
	_, _ = engine.Insert(user)

	// 插入多条记录
	//users := make([]User, 10)
	var users []User
	for i := range [5]int{} {
		fmt.Println(i)
		user := User{
			//Id:   int64(i),
			Name: fmt.Sprintf("aaa%d", i),
		}
		users = append(users, user)
	}
	result, err := engine.Insert(users)
	fmt.Printf("xxxx%v  -- %v  %v\n", users, result, err)

	// 执行sql查询
	sql := "select * from user"
	results, _ := engine.Query(sql)

	fmt.Println(results)

	// 执行sql命令
	sql = "update user set usr_name = 'aaaaa' where id % 2 = ?"
	res, err := engine.Exec(sql, 0)
	fmt.Println(res, err)
}

// 开启事务
func MyTransactionOps() error {
	session := engine.NewSession()
	defer session.Close()

	// // add Begin() before any action
	// if err := session.Begin(); err != nil {
	// 	return err
	// }

	// user1 := Userinfo{Username: "xiaoxiao", Departname: "dev", Alias: "lunny", Created: time.Now()}
	// if _, err := session.Insert(&user1); err != nil {
	// 	return err
	// }
	// user2 := Userinfo{Username: "yyy"}
	// if _, err = session.Where("id = ?", 2).Update(&user2); err != nil {
	// 	return err
	// }

	// if _, err = session.Exec("delete from userinfo where username = ?", user2.Username); err != nil {
	// 	return err
	// }

	// // add Commit() after all actions
	return session.Commit()
}
