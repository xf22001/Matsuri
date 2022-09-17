#!/bin/bash

#================================================================
#   
#   
#   文件名称：compile.sh
#   创 建 者：肖飞
#   创建日期：2022年08月20日 星期六 12时56分41秒
#   修改日期：2022年09月17日 星期六 16时45分35秒
#   描    述：
#
#================================================================
function main() {
	./run init action go && ./run lib core
	./run init action library
	#./gradlew tasks --all
	#./gradlew app:assembleFdroid --info
	#./gradlew app:installFdroidDebug --info
	#./gradlew app:assembleRelease --info
	./gradlew app:assemblePlay --info
}

main $@
