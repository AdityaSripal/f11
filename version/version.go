//nolint
package version

const Major = "0"
const Minor = "1"

// This will be overwritten during deployment by CI.
var Release = "0-dev"

var Version = Major + "." + Minor + "." + Release
