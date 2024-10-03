package client

import (
	"encoding/json"
	"fmt"
	sender "hacss/src/communication/clientsender"
	"hacss/src/config"
	logging "hacss/src/logging"
	"hacss/src/message"
	pb "hacss/src/proto/proto/communication"
	"hacss/src/utils"
	"log"
)

var cid int64
var err error
var clientTimer int

func GetCID() int64 {
	return cid
}

func SignedRequest(cid int64, dataSer []byte) ([]byte, bool) {
	request := message.MessageWithSignature{
		Msg: dataSer,
		Sig: []byte(""), //cryptolib.GenSig(cid, dataSer),
	}

	requestSer, err := request.Serialize()
	if err != nil {
		p := fmt.Sprintf("[Client error] fail to serialize the request with signiture: %v", err)
		logging.PrintLog(true, logging.ErrorLog, p)
		return []byte(""), false
	}

	return requestSer, true
}

func SendWriteRequest(op []byte) {

	dataSer, result1 := CreateRequest(cid, op)
	if !result1 {
		return
	}

	requestSer, result2 := SignedRequest(cid, dataSer)
	if !result2 {
		return
	}
	log.Println("len of request: ", len(requestSer))
	sender.BroadcastRequest(pb.MessageType_WRITE, requestSer)
}

func SendReconstructRequest(instance int) {

	op := utils.IntToBytes(instance)
	dataSer, result1 := CreateReconstructRequest(cid, op)
	if !result1 {
		return
	}

	requestSer, result2 := SignedRequest(cid, dataSer)
	if !result2 {
		return
	}
	log.Println("len of request: ", len(requestSer))
	sender.BroadcastRequest(pb.MessageType_RECONSTRUCT, requestSer)
}

func SendTestHacssRequest() {

	dataSer, result1 := CreateTestHacssRequest(cid)
	if !result1 {
		return
	}

	requestSer, result2 := SignedRequest(cid, dataSer)
	if !result2 {
		return
	}
	log.Println("len of request: ", len(requestSer))
	sender.BroadcastRequest(pb.MessageType_TEST_HACSS, requestSer)
}

func SendBatchRequest(op []byte, bitchSize int) {
	var requestArr [][]byte
	for i := 0; i < bitchSize; i++ {

		dataSer, result1 := CreateRequest(cid, op)
		if !result1 {
			return
		}

		requestSer, result2 := SignedRequest(cid, dataSer)
		if !result2 {
			return
		}
		//log.Println("len of request in batch: ",len(requestSer))
		requestArr = append(requestArr, requestSer)
	}
	byteRequsets, err := SerializeRequests(requestArr)
	if err != nil {
		log.Fatal("[Client error] fail to serialize the message.")
	}

	sender.BroadcastRequest(pb.MessageType_WRITE_BATCH, byteRequsets)
}

func CreateRequest(cid int64, op []byte) ([]byte, bool) {
	data := message.ClientRequest{
		Type: pb.MessageType_WRITE,
		ID:   cid,
		OP:   op,
		TS:   utils.MakeTimestamp(),
	}

	dataSer, err := data.Serialize()
	if err != nil {
		p := fmt.Sprintf("[Client error] fail to serialize the write request: %v", err)
		logging.PrintLog(true, logging.ErrorLog, p)
		return []byte(""), false
	}

	return dataSer, true
}

func CreateReconstructRequest(cid int64, instance []byte) ([]byte, bool) {
	data := message.ClientRequest{
		Type: pb.MessageType_RECONSTRUCT,
		ID:   cid,
		OP:   instance,
		TS:   utils.MakeTimestamp(),
	}

	dataSer, err := data.Serialize()
	if err != nil {
		p := fmt.Sprintf("[Client error] fail to serialize the write request: %v", err)
		logging.PrintLog(true, logging.ErrorLog, p)
		return []byte(""), false
	}

	return dataSer, true
}

func CreateTestHacssRequest(cid int64) ([]byte, bool) {
	data := message.ClientRequest{
		Type: pb.MessageType_TEST_HACSS,
		ID:   cid,
		TS:   utils.MakeTimestamp(),
	}

	dataSer, err := data.Serialize()
	if err != nil {
		p := fmt.Sprintf("[Client error] fail to serialize the write request: %v", err)
		logging.PrintLog(true, logging.ErrorLog, p)
		return []byte(""), false
	}

	return dataSer, true
}

/*
Serialize data into a json object in bytes
Output

	[]byte: serialized request
	error: err is nil if request is serialized
*/
func SerializeRequests(r [][]byte) ([]byte, error) {
	jsons, err := json.Marshal(r)
	if err != nil {
		p := fmt.Sprintf("[Client error] fail to serialize the message %v", err)
		logging.PrintLog(true, logging.ErrorLog, p)
		return []byte(""), err
	}
	return jsons, nil
}

func StartClient(rid string, loadkey bool) {
	logging.SetID(rid)
	config.LoadConfig()
	logging.SetLogOpt(config.FetchLogOpt())

	log.Printf("Client %s started.", rid)
	cid, err = utils.StringToInt64(rid)
	sender.StartClientSender(rid, loadkey)
	clientTimer = config.FetchBroadcastTimer()
}
