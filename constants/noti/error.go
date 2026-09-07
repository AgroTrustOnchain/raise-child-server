package noti

const (
	INTERNALL_ERR_MSG string = "There is something wrong in the system during the process. Please try again later."

	REPO_ERR_MSG string = "Error in %s repository at "

	MAIL_ERR_MSG string = "Error while generating mail at %s service -  "

	GIN_ERR_MSG string = "Error while starting gin server in %s service - "

	NET_LISTENING_ERR_MSG string = "Error while listening on port %s - "
)

// Env
const (
	ENV_LOAD_ERR_MSG string = "Error while loading .env - "

	ENV_SET_ERR_MSG string = "Error while setting environment variable %s with value %s - "
)

// Database
const (
	DB_CONNECTION_ERR_MSG string = "Error while connecting to database - "

	DB_MIGRATION_ERR_MSG string = "Error while migrating database in %s service - "

	DB_SET_CONNECTION_STRING_ERR_MSG string = "Error while setting database connection string in %s service - "
)

// On-chain
const (
	BUILDING_TX_ERR_MSG                      string = "Error while building transaction - "
	BATCHING_TX_ERR_MSG                      string = "Error while batching transaction - "
	EXECUTING_TX_ERR_MSG                     string = "Error while executing transaction - "
	SPLIT_COIN_ERR_MSG                       string = "Error while splitting coin - "
	GET_BALANCE_ERR_MSG                      string = "Error while getting balance - "
	GET_FAUCET_HOST_ERR_MSG                  string = "Error while getting faucet host - "
	FAUCET_ERR_MSG                           string = "Error while fauceting sui token - "
	RETRIEVE_OWNED_COINS_ERR_MSG             string = "Error while retrieving owned coins - "
	RETRIEVE_ON_CHAIN_DATA_ERR_MSG           string = "Error while retrieving on-chain data - "
	ApproRETRIEVE_DYNAMIC_FIELDS_ERR_MSGvers string = "Error while retrieving dynamic fields - "
)

// Payment
const (
	PAYMENT_INIT_ENV_ERR_MSG                 string = "Error while setup %s environment - "
	PAYMENT_GENERATE_TRANSACTION_URL_ERR_MSG string = "Error while generating %s transaction URL - "
)
