.PHONY: mytts
mytts: protogen
	bash install.sh

.PHONY: run
run: protogen
	@echo "Running MyTTS backend..."
	bash run.sh
	@echo "MyTTS backend started."

.PHONY: test
test: protogen
	@echo "Testing MyTTS backend..."
	bash test.sh
	@echo "MyTTS backend tested."

.PHONY: protogen
protogen: backend_pb2_grpc.py backend_pb2.py

.PHONY: protogen-clean
protogen-clean:
	$(RM) backend_pb2_grpc.py backend_pb2.py

backend_pb2.py backend_pb2_grpc.py: ../../backend.proto
	. venv/bin/activate && python3 -m grpc_tools.protoc -I../.. --python_out=. --grpc_python_out=. ../../backend.proto

.PHONY: clean
clean: protogen-clean
	rm -rf venv __pycache__ *.wav *.pyc
	find . -name "*.pyc" -delete
	find . -name "__pycache__" -type d -exec rm -rf {} +

.PHONY: install
install:
	bash install.sh
