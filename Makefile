all:
	go fmt
	go build

linux:
	$(MAKE) GOOS=linux all

TRASH=makeconst.cpp const.go
clean:
	rm $(TRASH)  || del $(TRASH)

test-dos:
	go run importconst.go -d -p dos \
	    FILE_ATTRIBUTE_NORMAL \
	    FILE_ATTRIBUTE_REPARSE_POINT \
	    FILE_ATTRIBUTE_HIDDEN \
	    CP_THREAD_ACP \
	    MOVEFILE_REPLACE_EXISTING \
	    MOVEFILE_COPY_ALLOWED \
	    MOVEFILE_WRITE_THROUGH

test-conio:
	go run importconst.go \
	    -d -p dos conio.h \
	    CTRL_CLOSE_EVENT \
	    CTRL_LOGOFF_EVENT \
	    CTRL_SHUTDOWN_EVENT \
	    CTRL_C_EVENT \
	    ENABLE_ECHO_INPUT \
	    ENABLE_PROCESSED_INPUT \
	    STD_INPUT_HANDLE \
	    STD_OUTPUT_HANDLE

test:
	go run importconst.go -d sample.h AHAHA IHIHI UFUFU
