package main

import (
	"fmt"

	"git.mmeiblog.cn/mei/aiComplain/internal"
)

func main() {
	// 先写点数据测试
	response, err := internal.SendAndAudit("2024级11班某女同学(🐭老鼠)于上周周五晚上在女生宿舍四处投毒，毒物为霉菌饼子，在霉菌里发现少量饼子。食用过的某同学说因为宿舍灯光暗，还以为是葱花饼。\r\n据说219(12班)宿舍三人受害严重，一人第二夜深夜直接被送回家，一人第三天下午回家，证物已上交老师，但该名女同学并无悔过之心，四处与他人就此事讲笑话，并单方面说受害同学原谅了她，事后还去各个宿舍索要食物。\r\n希望大家多加小心该名女同学，其外号在学校叫为老鼠🐭，可询问周围同学其真名，目前在2024级的11班")
	fmt.Println(response.Data, err)
}
