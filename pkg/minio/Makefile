test.thanosbench:
	rm -rf testdata
	mkdir -p testdata
	docker run --rm \
		--volume $(PWD):/thanosbench \
		--volume $(PWD)/resources:/etc/thanosbench \
		--workdir /thanosbench \
		quay.io/thanos/thanosbench:v0.3.0-rc.0 \
		block gen --output.dir=testdata --config-file=/etc/thanosbench/thanosbench.yaml