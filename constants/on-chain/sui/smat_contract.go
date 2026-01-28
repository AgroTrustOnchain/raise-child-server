package sui

// Modules
const (
	MODULE_MANAGE  string = "manage"
	MODULE_RECORD  string = "record"
	MODULE_CHILD   string = "child"
	MODULE_STAFF   string = "staff"
	MODULE_SPONSOR string = "sponsor"
	MODULE_POOL    string = "pool"
	MODULE_GIFT    string = "gift"
)

// Functions
const (
	DONATE_SUI_POOL_FUNCTION             string = "donate_to_sui_pool"
	WITHDRAW_FROM_SUI_POOL_FUNCTION      string = "withdraw_from_sui_pool"
	ADD_CHILD_FUNCTION                   string = "add_child"
	ADD_STRING_METADATA_FUNCTION         string = "add_string_metadata"
	ADD_NUMBER_METADATA_FUNCTION         string = "add_u64_metadata"
	UPDATE_STRING_METADATA_FUNCTION      string = "update_string_metadata"
	UPDATE_NUMBER_METADATA_FUNCTION      string = "update_u64_metadata"
	REGISTER_STAFF_FUNCTION              string = "register_staff"
	REGISTER_VOLUNTEER_FUNCTION          string = "register_volunteer"
	REGISTER_LOCAL_LEADER_FUNCTION       string = "register_local_leader"
	REGISTER_ADMIN_FUNCTION              string = "register_admin"
	DONATE_TO_POOL_FUNCTION              string = "donate_to_pool"
	DONATE_TO_LOCAL_POOL_FUNCTION        string = "donate_to_local_pool"
	WITHDRAW_FROM_POOL_FUNCTION          string = "withdraw_from_pool"
	CREATE_WITHDRAW_PROPOSAL_FUNCTION    string = "create_withdraw_proposal"
	VOTE_WITHDRAW_PROPOSAL_FUNCTION      string = "vote_withdraw_proposal"
	UPDATE_PUBLISHER_NFT_FUNCTION        string = "update_publisher_nft"
	CREATE_CENTER_FUNCTION               string = "create_children_center"
	MINT_REGISTER_VOLUNTEER_CAP_FUNCTION string = "mint_register_volunteer_cap"
	MINT_REGISTER_LEADER_CAP_FUNCTION    string = "mint_register_local_leader_cap"
	MINT_REGISTER_ADMIN_CAP_FUNCTION     string = "mint_register_admin_cap"
	MINT_UPLOAD_CENTER_CAP_FUNCTION      string = "mint_upload_center_cap"
	CREATE_GIFT_FUNCTION                 string = "create_gift"
	CANCEL_GIFT_FUNCTION                 string = "cancel_gift"
	CONFIRM_RECIEVE_GIFT_FUNCTION        string = "confirm_recieved"
)

// Structs
const (
	MANAGE_STRUCT             string = "Manage"
	SUI_POOL_STRUCT           string = "SuiPool"
	TRANSACTION_RECORD_STRUCT string = "TransactionRecord"
	CHILD_STRUCT              string = "Child"
	ADMIN_NFT_STRUCT          string = "AdminNFT"
	STAFF_NFT_STRUCT          string = "StaffNFT"
	SPONSOR_NFT_STRUCT        string = "SponsorNFT"
	GIFT_STRUCT               string = "Gift"
)

// Cap structs
const (
	ADMIN_CAP_STRUCT              string = "AdminCap"
	UPDATE_ADMIN_INFO_CAP_STRUCT  string = "UpdateAdminInfoAfterPublishCap"
	REGISTER_VOLUNTEER_CAP_STRUCT string = "RegisterVolunteerCap"
	REGISTER_LEADER_CAP_STRUCT    string = "RegisterLocalLeaderCap"
	REGISTER_ADMIN_CAP_STRUCT     string = "RegisterAdminCap"
	UPLOAD_CENTER_CAP_STRUCT      string = "UploadCenterCap"
)

// Events
const (
	TRANSACTION_RECORD_EVENT string = "TransactionRecordEvent"
)
