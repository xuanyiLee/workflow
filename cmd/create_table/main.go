package main

import (
	"datax/conf"
	"datax/pkg"
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

var mysqlCnf conf.MysqlConf
var moCnf conf.MysqlConf

func main() {
	cfg := conf.NewConf("./conf/matrixone.ini")
	err := cfg.Load()
	if err != nil {
		panic(err)
	}
	mysqlCnf = conf.MyCnf
	moCnf = conf.MoConf

	mysqlCnf.Password, moCnf.Password = pkg.GetPwd()
	switch cfg.Type {
	case "ssb":
		create_table("./ssb/bin/create-ssb-tables.sh")
	case "tpch":
		create_table("./tpch/bin/create-tpch-tables.sh")
	case "tpcc":
		cfg.CreatetpccConf()
		create_table("./tpcc/bin/create-tpcc-tables.sh")
	case "tpcds":
		create_table("./tpcds/bin/create-tpcds-tables.sh")
	}

	os.Exit(0)
}

func create_table(command string) {
	cmd := exec.Command(command)

	cmd.Env = append(os.Environ(), "HOST="+mysqlCnf.HOST, "PORT="+strconv.Itoa(mysqlCnf.Port), "USER="+mysqlCnf.Username, "PASSWORD="+mysqlCnf.Password, "DB="+mysqlCnf.DataBase)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(fmt.Sprintf("create mysql table fail: %v", err))
		os.Exit(1)
	}
	fmt.Println(string(output))

	cmd = exec.Command(command)
	cmd.Env = append(os.Environ(), "HOST="+moCnf.HOST, "PORT="+strconv.Itoa(moCnf.Port), "USER="+moCnf.Username, "PASSWORD="+moCnf.Password, "DB="+moCnf.DataBase)
	output, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Println(fmt.Sprintf("create matrixone table fail: %v", err))
		os.Exit(1)
	}

	fmt.Println(string(output))
}
