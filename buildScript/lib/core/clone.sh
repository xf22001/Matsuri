COMMIT=$(cat libcore/core_commit.txt)

cd ..
[ -d v2ray-core ] && exit 0
rm -rf v2ray-core
git clone --no-checkout https://github.com/xf22001/MatsuriDayo-v2ray-core.git v2ray-core
cd v2ray-core
git checkout $COMMIT
