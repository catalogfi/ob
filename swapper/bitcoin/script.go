package bitcoin

import (
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/txscript"
)

// OP_IF
//
//	    OP_SHA256
//	    ${secretHash}
//	    OP_EQUALVERIFY
//	    OP_DUP
//	    OP_HASH160
//	    ${redeemerAddress}
//	    OP_EQUALVERIFY
//	    OP_CHECKSIG
//	OP_ELSE
//	    ${waitTime}
//	    OP_CHECKSEQUENCEVERIFY
//	    OP_DROP
//	    OP_DUP
//	    OP_HASH160
//	    ${initiatorAddress}
//	    OP_EQUALVERIFY
//	    OP_CHECKSIG
//	OP_ENDIF
func NewHTLCScript(initiatorAddress, redeemerAddress btcutil.Address, secretHash []byte, waitTime int64) ([]byte, error) {
	return txscript.NewScriptBuilder().
		AddOp(txscript.OP_IF).
		AddOp(txscript.OP_SHA256).
		AddData(secretHash).
		AddOp(txscript.OP_EQUALVERIFY).
		AddOp(txscript.OP_DUP).
		AddOp(txscript.OP_HASH160).
		AddData(redeemerAddress.ScriptAddress()).
		AddOp(txscript.OP_EQUALVERIFY).
		AddOp(txscript.OP_CHECKSIG).
		AddOp(txscript.OP_ELSE).
		AddInt64(waitTime).
		AddOp(txscript.OP_CHECKSEQUENCEVERIFY).
		AddOp(txscript.OP_DROP).
		AddOp(txscript.OP_DUP).
		AddOp(txscript.OP_HASH160).
		AddData(initiatorAddress.ScriptAddress()).
		AddOp(txscript.OP_EQUALVERIFY).
		AddOp(txscript.OP_CHECKSIG).
		AddOp(txscript.OP_ENDIF).
		Script()
}

func NewHTLCRedeemScript(pubKey, secret []byte) ([]byte, error) {
	return txscript.NewScriptBuilder().
		AddData(pubKey).
		AddData(secret).
		AddOp(txscript.OP_TRUE).
		Script()
}

func NewHTLCRefundScript(pubKey []byte) ([]byte, error) {
	return txscript.NewScriptBuilder().
		AddData(pubKey).
		AddOp(txscript.OP_FALSE).
		Script()
}
