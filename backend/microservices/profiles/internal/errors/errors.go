package errors

import "github.com/Archisman-Mridha/x-clone/backend/pkg/utils"

var ErrDuplicateProfileID = utils.NewAPIError("profile with given id already exists")
