package defaults

//
// This file contains the input parameters for the lambda function.
// Lambda functions cannot take command-line parameters, hence they are baked into the binary.
// Alternatively it could be implemented to read them from DynamoDB or other storage.
//

const Major = "0"
const Minor = "1"
const ContentType = "application/json; charset=utf8"

// These will be overwritten during deployment by CI.
var Release = "0-dev"
var TestnetName = "devnet"
var FromKey = ""
var Amount = ""
var Node = ""

// Calculated value
var Version = Major + "." + Minor + "." + Release
