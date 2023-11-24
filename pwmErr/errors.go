package pwmErr

import "errors"

var InvalidOpt = errors.New("invalid option value")
var FailToCreateRepository = errors.New("fail to create some of repository")
var InsertDB = errors.New("fail to insert data")
var DBInit = errors.New("fail to initialize password storage")
var NoUser = errors.New("not found user")
var GeneratePw = errors.New("fail to generate hashed password")
var Unknown = errors.New("unknown err")
var NotRegisteredKey = errors.New("not registered key")
var ExitErr = errors.New("it is not a real error. use for exit")
var FailUpdatePw = errors.New("fail to update password")
var FailUpdatePwDescription = errors.New("fail to update password description")
var NoRecord = errors.New("no record")
var ConvertCsv = errors.New("fail to convert csv")
var WrongPw = errors.New("wrong password")
var AppVersionUnknown = errors.New("unknown err related to app version")
var OptParseErr = errors.New("fail to parse options")
var AlreadyExistKey = errors.New("already exist key")
