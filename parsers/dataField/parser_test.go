package datafield

import (
	"encoding/hex"
	"testing"

	"github.com/Dharitri-org/drtg-core/core"
	"github.com/Dharitri-org/drtg-vm-common/mock"
	"github.com/stretchr/testify/require"
)

func createMockArgumentsOperationParser() *ArgsOperationDataFieldParser {
	return &ArgsOperationDataFieldParser{
		Marshalizer:   &mock.MarshalizerMock{},
		AddressLength: 32,
	}
}

func TestNewOperationDataFieldParser(t *testing.T) {
	t.Parallel()

	t.Run("NilMarshalizer", func(t *testing.T) {
		t.Parallel()

		arguments := createMockArgumentsOperationParser()
		arguments.Marshalizer = nil

		_, err := NewOperationDataFieldParser(arguments)
		require.Equal(t, core.ErrNilMarshalizer, err)
	})

	t.Run("ShouldWork", func(t *testing.T) {
		t.Parallel()

		arguments := createMockArgumentsOperationParser()

		parser, err := NewOperationDataFieldParser(arguments)
		require.NotNil(t, parser)
		require.Nil(t, err)
	})
}

func TestParseQuantityOperationsDCT(t *testing.T) {
	t.Parallel()

	arguments := createMockArgumentsOperationParser()
	parser, _ := NewOperationDataFieldParser(arguments)

	t.Run("DCTLocalBurn", func(t *testing.T) {
		t.Parallel()

		dataField := []byte("DCTLocalBurn@4d4949552d616263646566@0102")
		res := parser.Parse(dataField, sender, sender, 3)
		require.Equal(t, &ResponseParseData{
			Operation: "DCTLocalBurn",
			DCTValues: []string{"258"},
			Tokens:    []string{"MIIU-abcdef"},
		}, res)
	})

	t.Run("DCTLocalMint", func(t *testing.T) {
		t.Parallel()

		dataField := []byte("DCTLocalMint@4d4949552d616263646566@1122")
		res := parser.Parse(dataField, sender, sender, 3)
		require.Equal(t, &ResponseParseData{
			Operation: "DCTLocalMint",
			DCTValues: []string{"4386"},
			Tokens:    []string{"MIIU-abcdef"},
		}, res)
	})

	t.Run("DCTLocalMintNotEnoughArguments", func(t *testing.T) {
		t.Parallel()

		dataField := []byte("DCTLocalMint@4d4949552d616263646566")
		res := parser.Parse(dataField, sender, sender, 3)
		require.Equal(t, &ResponseParseData{
			Operation: "DCTLocalMint",
		}, res)
	})
}

func TestParseQuantityOperationsNFT(t *testing.T) {
	t.Parallel()

	arguments := createMockArgumentsOperationParser()
	parser, _ := NewOperationDataFieldParser(arguments)

	t.Run("DCTNFTCreate", func(t *testing.T) {
		t.Parallel()

		dataField := []byte("DCTNFTCreate@4E46542D316630666638@01@4E46542D31323334@03e8@516d664132487465726e674d6242655467506b3261327a6f4d357965616f33456f61373678513775346d63646947@746167733a746573742c667265652c66756e3b6d657461646174613a5468697320697320612074657374206465736372697074696f6e20666f7220616e20617765736f6d65206e6674@0101")
		res := parser.Parse(dataField, sender, sender, 3)
		require.Equal(t, &ResponseParseData{
			Operation: "DCTNFTCreate",
			DCTValues: []string{"1"},
			Tokens:    []string{"NFT-1f0ff8"},
		}, res)
	})

	t.Run("DCTNFTBurn", func(t *testing.T) {
		t.Parallel()

		dataField := []byte("DCTNFTBurn@5454545454@0102@123456")
		res := parser.Parse(dataField, sender, sender, 3)
		require.Equal(t, &ResponseParseData{
			Operation: "DCTNFTBurn",
			DCTValues: []string{"1193046"},
			Tokens:    []string{"TTTTT-0102"},
		}, res)
	})

	t.Run("DCTNFTAddQuantity", func(t *testing.T) {
		t.Parallel()

		dataField := []byte("DCTNFTAddQuantity@5454545454@02@03")
		res := parser.Parse(dataField, sender, sender, 3)
		require.Equal(t, &ResponseParseData{
			Operation: "DCTNFTAddQuantity",
			DCTValues: []string{"3"},
			Tokens:    []string{"TTTTT-02"},
		}, res)
	})

	t.Run("DCTNFTAddQuantityNotEnoughArguments", func(t *testing.T) {
		t.Parallel()

		dataField := []byte("DCTNFTAddQuantity@54494b4954414b41@02")
		res := parser.Parse(dataField, sender, sender, 3)
		require.Equal(t, &ResponseParseData{
			Operation: "DCTNFTAddQuantity",
		}, res)
	})
}

func TestParseBlockingOperationDCT(t *testing.T) {
	t.Parallel()

	arguments := createMockArgumentsOperationParser()
	parser, _ := NewOperationDataFieldParser(arguments)

	t.Run("DCTFreeze", func(t *testing.T) {
		t.Parallel()

		dataField := []byte("DCTFreeze@5454545454")
		res := parser.Parse(dataField, sender, receiver, 3)
		require.Equal(t, &ResponseParseData{
			Operation: "DCTFreeze",
			Tokens:    []string{"TTTTT"},
		}, res)
	})

	t.Run("DCTFreezeNFT", func(t *testing.T) {
		t.Parallel()

		dataField := []byte("DCTFreeze@544f4b454e2d616263642d3031")
		res := parser.Parse(dataField, sender, receiver, 3)
		require.Equal(t, &ResponseParseData{
			Operation: "DCTFreeze",
			Tokens:    []string{"TOKEN-abcd-01"},
		}, res)
	})

	t.Run("DCTWipe", func(t *testing.T) {
		t.Parallel()

		dataField := []byte("DCTWipe@534b4537592d37336262636404")
		res := parser.Parse(dataField, sender, receiver, 3)
		require.Equal(t, &ResponseParseData{
			Operation: "DCTWipe",
			Tokens:    []string{"SKE7Y-73bbcd-04"},
		}, res)
	})

	t.Run("DCTFreezeNoArguments", func(t *testing.T) {
		t.Parallel()

		dataField := []byte("DCTFreeze")
		res := parser.Parse(dataField, sender, receiver, 3)
		require.Equal(t, &ResponseParseData{
			Operation: "DCTFreeze",
		}, res)
	})

	t.Run("SCCall", func(t *testing.T) {
		t.Parallel()

		dataField := []byte("callMe@01")
		res := parser.Parse(dataField, sender, receiverSC, 3)
		require.Equal(t, &ResponseParseData{
			Operation: operationTransfer,
			Function:  "callMe",
		}, res)
	})
}

func TestOperationDataFieldParser_ParseRelayed(t *testing.T) {
	t.Parallel()

	args := createMockArgumentsOperationParser()
	parser, _ := NewOperationDataFieldParser(args)

	t.Run("RelayedTxOk", func(t *testing.T) {
		t.Parallel()

		dataField := []byte("relayedTx@7b226e6f6e6365223a362c2276616c7565223a302c227265636569766572223a2241414141414141414141414641436e626331733351534939726e6d697a69684d7a3631665539446a71786b3d222c2273656e646572223a2248714b386459464a43474144346a756d4e4e742b314530745a6579736376714c7a38624c47574e774177453d222c226761735072696365223a313030303030303030302c226761734c696d6974223a31353030303030302c2264617461223a2252454e5556484a68626e4e6d5a584a414e444d304e7a526a4e4451795a444d354d7a497a4f444d304d7a6b7a4d6b41774d325534514459794e7a55334f54517a4e6a67324e54637a4e7a5241595441774d4441774d44413d222c22636861696e4944223a2252413d3d222c2276657273696f6e223a312c227369676e6174757265223a2262367331755349396f6d4b63514448344337624f534a632f62343166577a3961584d777334526966552b71343870486d315430636f72744b727443484a4258724f67536b3651333254546f7a6e4e2b7074324f4644413d3d227d")

		res := parser.Parse(dataField, sender, receiver, 3)

		rcv, _ := hex.DecodeString("0000000000000000050029db735b3741223dae79a2ce284ccfad5f53d0e3ab19")
		require.Equal(t, &ResponseParseData{
			IsRelayed:        true,
			Operation:        "DCTTransfer",
			Function:         "buyChest",
			Tokens:           []string{"CGLD-928492"},
			DCTValues:        []string{"1000"},
			Receivers:        [][]byte{rcv},
			ReceiversShardID: []uint32{1},
		}, res)
	})

	t.Run("RelayedTxV2ShouldWork", func(t *testing.T) {
		t.Parallel()

		dataField := []byte(core.RelayedTransactionV2 +
			"@" +
			hex.EncodeToString(receiverSC) +
			"@" +
			"0A" +
			"@" +
			hex.EncodeToString([]byte("callMe@02")) +
			"@" +
			"01a2")

		res := parser.Parse(dataField, sender, receiver, 3)
		require.Equal(t, &ResponseParseData{
			IsRelayed:        true,
			Operation:        operationTransfer,
			Function:         "callMe",
			Receivers:        [][]byte{receiverSC},
			ReceiversShardID: []uint32{0},
		}, res)
	})

	t.Run("RelayedTxV2NotEnoughArgs", func(t *testing.T) {
		t.Parallel()

		dataField := []byte(core.RelayedTransactionV2 + "@abcd")
		res := parser.Parse(dataField, sender, receiver, 3)
		require.Equal(t, &ResponseParseData{
			IsRelayed: true,
		}, res)
	})

	t.Run("RelayedTxV1NoArguments", func(t *testing.T) {
		t.Parallel()

		dataField := []byte(core.RelayedTransaction)
		res := parser.Parse(dataField, sender, receiver, 3)
		require.Equal(t, &ResponseParseData{
			IsRelayed: true,
		}, res)
	})

	t.Run("RelayedTxV2WithRelayedTxIn", func(t *testing.T) {
		t.Parallel()

		dataField := []byte(core.RelayedTransactionV2 +
			"@" +
			hex.EncodeToString(receiverSC) +
			"@" +
			"0A" +
			"@" +
			hex.EncodeToString([]byte(core.RelayedTransaction)) +
			"@" +
			"01a2")
		res := parser.Parse(dataField, sender, receiver, 3)
		require.Equal(t, &ResponseParseData{
			IsRelayed: true,
		}, res)
	})

	t.Run("RelayedTxV2WithNFTTransfer", func(t *testing.T) {
		t.Parallel()

		nftTransferData := []byte("DCTNFTTransfer@4c4b4641524d2d396431656138@34ae14@728faa2c8883760aaf53bb@000000000000000005001e2a1428dd1e3a5146b3960d9e0f4a50369904ee5483@636c61696d5265776172647350726f7879@00000000000000000500a655b2b534218d6d8cfa1f219960be2f462e92565483")
		dataField := []byte(core.RelayedTransactionV2 +
			"@" +
			hex.EncodeToString(receiver) +
			"@" +
			"0A" +
			"@" +
			hex.EncodeToString(nftTransferData) +
			"@" +
			"01a2")
		res := parser.Parse(dataField, sender, receiver, 3)
		rcv, _ := hex.DecodeString("000000000000000005001e2a1428dd1e3a5146b3960d9e0f4a50369904ee5483")
		require.Equal(t, &ResponseParseData{
			IsRelayed:        true,
			Operation:        "DCTNFTTransfer",
			DCTValues:        []string{"138495980998569893315957691"},
			Tokens:           []string{"LKFARM-9d1ea8-34ae14"},
			Receivers:        [][]byte{rcv},
			ReceiversShardID: []uint32{1},
			Function:         "claimRewardsProxy",
		}, res)
	})

	t.Run("DCTTransferRole", func(t *testing.T) {
		t.Parallel()

		dataField := []byte("DCTNFTCreateRoleTransfer@01010101@020202")
		res := parser.Parse(dataField, sender, receiver, 3)
		require.Equal(t, &ResponseParseData{
			Operation: "DCTNFTCreateRoleTransfer",
		}, res)
	})
}

func TestParseSCDeploy(t *testing.T) {
	arguments := createMockArgumentsOperationParser()
	parser, _ := NewOperationDataFieldParser(arguments)

	t.Run("ScDeploy", func(t *testing.T) {
		t.Parallel()

		dataField := []byte("0101020304050607")
		rcvAddr := make([]byte, 32)

		res := parser.Parse(dataField, sender, rcvAddr, 3)
		require.Equal(t, &ResponseParseData{
			Operation: operationDeploy,
		}, res)
	})
}

func TestGuardians(t *testing.T) {
	arguments := createMockArgumentsOperationParser()
	parser, _ := NewOperationDataFieldParser(arguments)

	t.Run("SetGuardian", func(t *testing.T) {
		t.Parallel()

		dataField := []byte("SetGuardian")

		res := parser.Parse(dataField, sender, sender, 3)
		require.Equal(t, &ResponseParseData{
			Operation: core.BuiltInFunctionSetGuardian,
		}, res)
	})

	t.Run("GuardAccount", func(t *testing.T) {
		t.Parallel()

		dataField := []byte("GuardAccount")

		res := parser.Parse(dataField, sender, sender, 3)
		require.Equal(t, &ResponseParseData{
			Operation: core.BuiltInFunctionGuardAccount,
		}, res)
	})

	t.Run("UnGuardAccount", func(t *testing.T) {
		t.Parallel()

		dataField := []byte("UnGuardAccount")

		res := parser.Parse(dataField, sender, sender, 3)
		require.Equal(t, &ResponseParseData{
			Operation: core.BuiltInFunctionUnGuardAccount,
		}, res)
	})
}
