package httperrors

// Bad Request
const IDNegative = "id must be integer greater than 0"

const ParamIDInvalid = "id in params must be a number"

const NoField = "no fields in request body"

const NameNull = "name cannot be null or empty"

const PhoneInvalid = "phone Number Invalid"

const BloodInvalid = "blood Group must be one of ['+A', '-A', '+B', '-B', '+O', '-O', '+AB', '-AB']"

// Not Found
const PatientNotFound = "patient not found"

// Internal Server Errors
const PatientFailed = "failed to GET 'Patient'"

const CreateFailed = "failed to CREATE 'Patient'"

const UpdateFailed = "failed to UPDATE 'Patient'"

const DeleteFailed = "failed to DELETE 'Patient'"
