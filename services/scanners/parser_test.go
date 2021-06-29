package scanners

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"oasisTracker/dmodels/oasis"
	"oasisTracker/smodels/container"
	"os"
	"reflect"
	"testing"

	beaconAPI "github.com/oasisprotocol/oasis-core/go/beacon/api"
	"github.com/oasisprotocol/oasis-core/go/common/cbor"
	"github.com/oasisprotocol/oasis-core/go/common/grpc"
	"github.com/oasisprotocol/oasis-core/go/consensus/api/transaction"
	stakingAPI "github.com/oasisprotocol/oasis-core/go/staking/api"
	grpcCommon "google.golang.org/grpc"
)

func TestParser_TestBalances(t *testing.T) {
	grpcTarget := "127.0.0.1:9101"
	grpcConn, err := grpc.Dial(grpcTarget, grpcCommon.WithInsecure())
	if err != nil {
		log.Print(err)
		t.Error(err)
		return
	}
	defer grpcConn.Close()

	c := &ParseContainer{
		blocks:          container.NewBlocksContainer(),
		blockSignatures: container.NewBlockSignatureContainer(),
		txs:             container.NewTxsContainer(),
		balances:        container.NewAccountsContainer(),
		rewards:         container.NewRewardsContainer(),
	}

	parser, err := NewParserTask(context.Background(), grpcConn, 0, c)
	if err != nil {
		t.Error(err)
		return
	}

	adr := stakingAPI.Address{}
	err = adr.UnmarshalText([]byte("oasis1qqekv2ymgzmd8j2s2u7g0hhc7e77e654kvwqtjwm"))
	if err != nil {
		t.Error(err)
		return
	}

	height := int64(3816110)
	acc, err := parser.stakingAPI.Account(context.Background(), &stakingAPI.OwnerQuery{
		Height: height,
		Owner:  adr,
	})
	log.Print(height, " ", acc)

	acc, err = parser.stakingAPI.Account(context.Background(), &stakingAPI.OwnerQuery{
		Height: height + 1,
		Owner:  adr,
	})
	log.Print(height+1, " ", acc)

	acc, err = parser.stakingAPI.Account(context.Background(), &stakingAPI.OwnerQuery{
		Height: height + 2,
		Owner:  adr,
	})
	log.Print(height+2, " ", acc)
}

func TestMainNet(t *testing.T) {
	grpcTarget := "127.0.0.1:9101"
	grpcConn, err := grpc.Dial(grpcTarget, grpcCommon.WithInsecure())
	if err != nil {
		log.Print(err)
		t.Error(err)
		return
	}

	b := beaconAPI.NewBeaconClient(grpcConn)

	ch, close, err := b.WatchEpochs(context.Background())
	if err != nil {
		t.Error(err)
	}

	defer close.Close()

	select {
	case ep := <-ch:
		log.Print(ep)

	}
}

func TestT2(t *testing.T) {
	//base := uint64(5047)
	//log.Print(6356 - 5046)
	//log.Print((3814764 - 1308) / 600)
	//log.Print(3814764 / 600)
	//log.Print(3814764 % 600)
	//log.Print(base + 564)

	grpcTarget := "127.0.0.1:9101"
	grpcConn, err := grpc.Dial(grpcTarget, grpcCommon.WithInsecure())
	if err != nil {
		log.Print(err)
		t.Error(err)
		return
	}

	b := beaconAPI.NewBeaconClient(grpcConn)

	//log.Print(b.GetEpoch(context.Background(), 3398334))

	log.Print(b.GetEpochBlock(context.Background(), 5662))
	log.Print(b.GetEpochBlock(context.Background(), 5663))

	base, _ := b.GetBaseEpoch(context.Background())
	log.Print(base)
	//ep, _ := b.GetEpoch(context.Background(), 3398334)
}

func TestParser_ParseEscrow(t *testing.T) {
	grpcTarget := "127.0.0.1:9101"
	grpcConn, err := grpc.Dial(grpcTarget, grpcCommon.WithInsecure())
	if err != nil {
		log.Print(err)
		t.Error(err)
		return
	}
	defer grpcConn.Close()

	c := &ParseContainer{
		blocks:          container.NewBlocksContainer(),
		blockSignatures: container.NewBlockSignatureContainer(),
		txs:             container.NewTxsContainer(),
		balances:        container.NewAccountsContainer(),
		rewards:         container.NewRewardsContainer(),
	}

	parser, err := NewParserTask(context.Background(), grpcConn, 0, c)
	if err != nil {
		t.Error(err)
		return
	}

	adr := stakingAPI.Address{}
	err = adr.UnmarshalText([]byte("oasis1qqekv2ymgzmd8j2s2u7g0hhc7e77e654kvwqtjwm"))
	if err != nil {
		t.Error(err)
		return
	}

	accInfo, err := parser.stakingAPI.Account(parser.ctx, &stakingAPI.OwnerQuery{
		Height: 3027600,
		Owner:  adr,
	})
	if err != nil {
		t.Error(err)
		return
	}

	log.Print(accInfo)
	return

	baseEpoch, _ := parser.beaconAPI.GetBaseEpoch(context.Background())

	log.Print(baseEpoch)

	for i := uint64(6359); i < 6369; i++ {

		log.Print("I: ", i, " ", uint64(i-1)*600+(i-uint64(baseEpoch)))

	}
	//3027601

	//height := int64(3827956)
	//height := int64(
	//	3814764)

	bl, err := parser.beaconAPI.GetEpochBlock(context.Background(), 6359)

	log.Print(bl)

	change := int64(6359 - 5047)

	log.Print(bl / change)

	block, err := parser.consensusAPI.GetBlock(context.Background(), 3027601)
	if err != nil {
		t.Error(err)
		return
	}

	b := oasis.Block{}
	//Nil pointer err
	err = cbor.Unmarshal(block.Meta, &b)
	if err != nil {
		t.Error(err)
		return
	}

	log.Print(b.Header.Height)
	//log.Print(b.Header.)

	//txsWithResults, err := parser.consensusAPI.GetTransactionsWithResults(context.Background(),
	//	height)
	//if err != nil {
	//	t.Error(err)
	//	return
	//}
	//
	//for _, value := range txsWithResults.Results {
	//	for _, value := range value.Events {
	//		log.Print(value.Staking)
	//	}
	//
	//	//log.Print(value.Events[0].Staking.Transfer)
	//}
}

func TestParser_ParseBase(t *testing.T) {

	grpcTarget := "127.0.0.1:9101"
	grpcConn, err := grpc.Dial(grpcTarget, grpcCommon.WithInsecure())
	if err != nil {
		log.Print(err)
		t.Error(err)
		return
	}
	defer grpcConn.Close()

	c := &ParseContainer{
		blocks:          container.NewBlocksContainer(),
		blockSignatures: container.NewBlockSignatureContainer(),
		txs:             container.NewTxsContainer(),
		balances:        container.NewAccountsContainer(),
		rewards:         container.NewRewardsContainer(),
	}

	parser, err := NewParserTask(context.Background(), grpcConn, 0, c)
	if err != nil {
		t.Error(err)
		return
	}

	ctx := context.Background()

	txsWithResults, err := parser.consensusAPI.GetTransactionsWithResults(ctx, 4518525)
	if err != nil {
		t.Error(err)
		return
	}

	for _, value := range txsWithResults.Transactions {
		tx := transaction.SignedTransaction{}

		err = cbor.Unmarshal(value, &tx)
		if err != nil {
			t.Error(err)
			return
		}

		in, _ := tx.PrettyType()
		log.Printf("%+v", in)

		raw := transaction.Transaction{}
		err = cbor.Unmarshal(tx.Blob, &raw)
		if err != nil {
			t.Error(err)
			return
		}

		log.Print(raw.Method.BodyType())
	}

	//err = parser.ParseBase(4518525)
	//if err != nil {
	//	t.Errorf("ParseBase error: %s", err)
	//	return
	//}

	//Check txs

	//parser.stakingAPI.DebondingDelegations()

	//log.Print(hex.EncodeToString(base58.Decode("90F8bf6A479f320ead074411a4B0e7944Ea8c9C1")))
	//
	//return
	//3033998
	//3817136

	//err = parser.parseBlockTransactions(oasis.Block{Header: types.Header{
	//	Height: 916707,
	//}})
	//if err != nil {
	//	t.Error(err)
	//	return
	//}

	log.Print("Success")
}

func TestT(t *testing.T) {
	data, err := hex.DecodeString("a463666565a26367617319271066616d6f756e744064626f6479ac6176026269645820000000000000000000000000000000000000000000000000000000000000ff03646b696e64016767656e65736973a465726f756e640065737461746586825422aa096e896e16c0b99bdb93d6084c947c3bddab5847f8458080a0d7e69c4a8bc3d3359b906926117ee0e7d163a44bf3b7bf366c5cc7822e2535aaa0c5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a4708082583522aa096e896e16c0b99bdb93d6084c947c3bddab0200000000000000000000000000000000000000000000000000000000000000005821a0000000000000000000000000000000000000000000000000000000006041bad082583522aa096e896e16c0b99bdb93d6084c947c3bddab0200000000000000000000000000000000000000000000000000000000000000015821a0000000000000000000000000000000000000000000000000000000000000000082583522aa096e896e16c0b99bdb93d6084c947c3bddab0200000000000000000000000000000000000000000000000000000000000000025821a0000000000000000000000000000000000000000000000000000000000000000082583522aa096e896e16c0b99bdb93d6084c947c3bddab0200000000000000000000000000000000000000000000000000000000000000035821a00000000000000000000000000000000000000000000000010000000000000001825367656e657369735f696e697469616c697a656441016a73746174655f726f6f7458208c5caf263e6262f256af1fcc38d66d3aa274bdc78c250d4eb3883b7fbdec582a7073746f726167655f7265636569707473f66773746f72616765a76a67726f75705f73697a65026d6d61785f6170706c795f6f70730273636865636b706f696e745f696e74657276616c1903e873636865636b706f696e745f6e756d5f6b6570740575636865636b706f696e745f6368756e6b5f73697a651a01000000756d696e5f77726974655f7265706c69636174696f6e02781b6d61785f6170706c795f77726974655f6c6f675f656e74726965731a000186a0686578656375746f72a56a67726f75705f73697a65026c6d61785f6d6573736167657318206d726f756e645f74696d656f7574057167726f75705f6261636b75705f73697a650072616c6c6f7765645f7374726167676c657273006876657273696f6e73a16776657273696f6ea069656e746974795f696458204430b20114c8769886ef92fdbd55c0608851093a12adc85992fdc9516dcf8e2f6c7465655f6861726477617265006d74786e5f7363686564756c6572a569616c676f726974686d6673696d706c656e6d61785f62617463685f73697a651927107362617463685f666c7573685f74696d656f75741a3b9aca00746d61785f62617463685f73697a655f62797465731a010000007570726f706f73655f62617463685f74696d656f7574057061646d697373696f6e5f706f6c696379a170656e746974795f77686974656c697374a168656e746974696573b83458200416dad177c8e16493ae813efd819c16d8716c40b98d02615a19748c807091eda05820053009343c9debf510fa97e14c374f1aca2d87327e0b9fc2fe0e766b70c894cca058200bb85eec124fae056bb06466e26f8e15c1e570d60e10f48679eaf1bd5fc1b010a058200c86cf69c58612279a964fb6dba79811a7d898a2690fdf0f6285a9fed70f5c92a05820143a9198cd45ca16866acfa5aae5a61802a032c536ae3ed41120259ce1edc5caa0582017979db97e26e0efe8b7751c16cde2424d593d4b86e3792ec4623db19cfde6e5a0582018358ed3e47ba92ad60746ef4fe5762b3eaf947726e999da350336fc5b05a840a058202769f0957b9810f359d26307d8f860e51a1f6d9ce3eb9c43bd030d772f498b41a05820277a2ee66d5a6ebc1c6d8c1392c4a3d478f6d50a038f2ad2923078d2b10a5781a058202a6312f91ce343dbe6f49d9a5fd6228499e699273f201aae0b2b1b84fc550c53a058202da28950a9dd02a1f9ec3544c6ca538dd45344d189392a80af3334a66f11119ca058203ad9e73f11bd53557889e0d60aea81b2be01ce26009b0e88a450696b70a928f1a058204430b20114c8769886ef92fdbd55c0608851093a12adc85992fdc9516dcf8e2fa05820459a304d799f4fe6fa1fa56ecfded5102232d9cabcb72b538a5bfc0316edbcefa058204f993b3ed391d35a19add9d9bde0e93bd005a4c52111111993b59249f9bc1ad8a05820502d0069ad82cdcdb91b9f97d873a2a3a1f4c10b81251e79d1c84213692278dda058205055c2a5cbd704e1dbc6d39b1b8a6c19c9fe2e064e79dbc30d402a55e9e0a4f9a05820524c234b562f11f1f1f5be8c313e50d56bc26376969f6971443b01fcfc3ecc69a0582056cac73d270175ed4551bf810d638dbcbc71c0a38a7af2c1ec545b5006a217e6a0582056deb715f17334c8149e0accbe72a4b3e89cc928b92f09239e224424c6408a63a05820597b3b1250659b7d256b67c6e28643a6e6de16ee6c2a48c355e7e55a8f98b880a0582059acc8efc94c7268f2087e7ee512a491f393511f9785e1c8a2196a32ef9af40ba058205aa057f91c8fd52f16b02bb1ae4c24d4a504aa93f190c498caf974ae32fe342aa058205d240b068b9b84eabf5fceaf2888f85f3a2ea351617b7f7e90b17fd8f0e02a93a058206031d8cff47e63ba42a1d8668240aace8aaa379e20cd17d51397e366b897f912a0582063ebdef8ed21941728ff5633b85383227ec9a370f973d04ea9dba173b29933c1a0582068e79fc6de4c52e5c46c4af1a43f99b6d25687fdddf627feb7e3c14de0202a9aa058206bf2d515fa08dda4686ac66cead3ff2ef9f409b1ba464bc41581dc551fd5a3d3a058206d3a24d1e97c19b994cd300081c43bf2ec30fd3b20796c17a4cdcdd92e3da814a058206dbd735a9d20cd2d628f5398dcacf10b404935daa4022276eb3de20fcb25b7b7a058208b3cb7f00994794e1be0bdb7915cb8065c5b90718e3f2e0a3d248b9c9825ad5da058208f34d662291e37a348b8c43a9788aa743385f63c9f3891b615c9d068d067f7eba0582090dd35bfb2ffc2fea1beb7756a4f37a6009e7460e7a606c62559177b3854509fa0582091faf603a2ba4e5be1426e27cfcf07733ce477602ae4f964c52a67994501005ba0582092ea56dcfb745cc78445291d0d6c99314e2866b934c06cac557572a97debca57a05820a0f7a482121a866c8c533077d4815b5d077a0ca768ee59ce62c5c8a692b62165a05820b038bd6711d807eac74e9561e1a6cd1570cc45279c7c67b84336f218af198108a05820b5ca27dfcac1af774799446e5ce4d41690cf2b013878e53cd85d76d6b5fd7ac1a05820b6f3e59861efe5aa242fd30c06a3409123b282574922545399222f0b0ba1d2b5a05820bb14a4bc5bbac78308615f8cd55ac3bb79bfa9b003b3fd407b799601c1269da4a05820c5037a7df2d275ce7511f110d81ce5b4ad62598030e98d4290101b1739618444a05820db5fa23eefe899804ded7e5c518e109c3e1bf5556e0225bfbbcb8006ab765633a05820dc344c30cef27b61d112f93ccf3d3a3a3306e81b968ecf383cb75260284d1c0ca05820dcd28c4ff433d757b6c6497048257fcde6b02c4dc120efd19f82df695e3b4092a05820df6292e59f4b5930f650450eaa147067226f6f495bd768a213d518828163b0a4a05820e206ce394d3d6dccaf339dd4a759539cffac2dbd1f7a7889bb439c5143c2052ba05820e20da9dfa5a83b3747c97abab4e8c56f79efeb2f74ff3c9071bd5e41ccb9768aa05820ec3c0a082415f531cd2f24c42d04c109634f207649740abba869350fbbf1db3aa05820f7e94ce6f43583e715ce368d33bf20c70d7251e8158f9eed1c36e937b0ca22fca05820f8c060cc8f6697326c3b12035626fb5ec6385f86b862a5faedda47f08b7db17ea05820fc5a22e6465d05cef41cf9a582ebd8e80cc120b36b3623c380a40b3884ee2042a05820ff295675a89dd830e523804c5645fa80047a79a05894ba251db8c29878adaf37a070676f7665726e616e63655f6d6f64656c01656e6f6e63651850666d6574686f64781872656769737472792e526567697374657252756e74696d65")
	if err != nil {
		fmt.Println("err", err)
		os.Exit(1)
	}
	var obj transaction.Transaction

	if err := cbor.Unmarshal(data, &obj); err != nil {
		fmt.Println("err:", err)
		os.Exit(1)
	}

	log.Print(reflect.TypeOf(obj.Method.BodyType()))

	var obj2 oasis.UntrustedRawValue
	//obj2 := map[string]interface{}{}

	if err := cbor.Unmarshal(data, &obj2); err != nil {
		fmt.Println("err:", err)
		os.Exit(1)
	}
	log.Print(obj)

	obj.PrettyPrint(context.Background(), "", os.Stdout)
}
