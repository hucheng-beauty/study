package unique

import (
    "github.com/google/uuid"
)

func UUID() (string, error) {
    id, err := uuid.NewUUID()
    if err != nil {
        return "", err
    }
    return id.String(), nil
}

func GlobalID() (globalID string, err error) {
    return UUID()
}
