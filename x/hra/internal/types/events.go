package types

const (
	EventTypeRegister    		= "register_name"
	EventTypeRegisterV2		 	= "register_name_v2"
	EventTypeRenew 				= "renew_name"
	EventTypeSetPrice 			= "set_price"
	EventTypeDelete				= "delete_name"
	EventTypeBuy				= "buy_name"
	EventTypeTransferName		= "transfer_name"
	EventTypeRegisterAddress 	= "register_address"
	EventTypeRemoveAddress		= "remove_address"
	EventTypeExpiredName 		= "expired_name"
	EventTypeRegisterBlockchainId = "RegisterBlockchainId"
	EventTypeRemoveBlockchainId   = "RemoveBlockchainId"

	AttributeKeySender				= "sender"
	AttributeKeyName  				= "name"
	AttributeKeyExpires				= "expires"
	AttributeKeyNewOwner 			= "new_owner"
	AttributeKeyPrice 				= "price"
	AttributeKeyBlockchainAddress 	= "blockchain_address"
	AttributeKeyBlockchainId		= "blockchain_id"
	AttributeKeyIndex				= "index"
	AttributeKeyTitle				= "title"
	AttributeKeyDescription			= "description"
	AttributeKeyReferral			= "referral"

	AttributeValueModule = ModuleName
)
