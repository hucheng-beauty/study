package model

import (
    "database/sql/driver"
    "encoding/json"
    "fmt"
)

type HarvestLimitations []HarvestLimit

// HarvestLimit is conception within SaaS space
type HarvestLimit struct {
    // It's OK to use CropType isolated, because the system consume
    // the HarvestLimit will not care about the meta &
    // the Limit value are always normalized
    CropType CropType `json:"type"`

    CropTypeMeta *CropTypeMeta `json:"crop_type_with_meta,omitempty"`

    // Limit is the limit of CropType, it's value are always
    // normalized(durationUnit: ms, timesUnit: 1)
    Limit int64 `json:"limit"`
}

const (
    CropTypeAccess      CropType = "access"
    CropTypeQPS         CropType = "qps"
    CropTypeConcurrency CropType = "concurrency"
    CropTypeQPM         CropType = "qpm"
    CropTypeTPM         CropType = "tpm"

    // the next 4 CropType have two meaning:
    // If it's apply to api side & used as usage quota,
    // it means usage quota auto refresh every 1 day
    // If it's not apply to api side,e.g. as a volcano config,
    // it only means a charge unit, it is something which Lifetime CropType not have
    // !Note: 除了在API侧使用以外，这4种类型在一些时候也被用为一种普通的表述单位，
    // 即他们仅仅代表一种计量单位，而不在乎是 lifeTime 还是什么

    CropTypeAudioDuration              CropType = "audio_duration" // default refresh period 1 DAY
    CropTypeTextBytes                  CropType = "text_bytes"     // default refresh period 1 DAY
    CropTypeRequests                   CropType = "requests"       // not support in api's online traffic at the moment
    CropTypeSessions                   CropType = "sessions"       // default refresh period 1 DAY
    CropTypeTextWords                  CropType = "text_words"     // default refresh period 1 DAY
    CropTypeTokens                     CropType = "tokens"
    CropTypeInputTextTokens            CropType = "input_text_tokens"
    CropTypeInputAudioTokens           CropType = "input_audio_tokens"
    CropTypeCachedTextTokens           CropType = "cached_text_tokens"
    CropTypeCachedAudioTokens          CropType = "cached_audio_tokens"
    CropTypeOutputTextTokens           CropType = "output_text_tokens"
    CropTypeOutputAudioTokens          CropType = "output_audio_tokens"
    CropTypeAudioDurationSingleFeature CropType = "audio_duration_single_feature"
    CropTypeAudioDurationAllFeatures   CropType = "audio_duration_all_features"
    CropTypeLifetimeAudioDuration      CropType = "audio_duration_lifetime"
    CropTypeLifetimeRequestCount       CropType = "requests_lifetime"
    CropTypeLifetimeSessionCount       CropType = "sessions_lifetime"
    CropTypeLifetimeTextWords          CropType = "text_words_lifetime"
    CropTypeLifetimeTokens             CropType = "tokens_lifetime"
)

// CropType marked Lifetime expected to never expire & never auto reset
type CropType string

var (
    // for the **Collection Below:
    // it'll be easy to use switch { case lo.Contains(**Collection, targetCropType): }
    // to check if a CropType Belongs to certain collection
    // but be careful: it's only a good idea to use exclusive CropType Collections below

    // QuotaTypes & UsageTypes exclusive to each other

    QuotaTypes = []CropType{
        CropTypeQPS, CropTypeConcurrency,
        CropTypeQPM, CropTypeTPM,
    }
    UsageTypes = []CropType{
        CropTypeAudioDuration, CropTypeTextBytes,
        CropTypeRequests, CropTypeSessions,
        CropTypeLifetimeAudioDuration, CropTypeLifetimeRequestCount,
        CropTypeLifetimeSessionCount, CropTypeTextWords,
        CropTypeLifetimeTextWords, CropTypeLifetimeTokens,
    }
    TokenTypes = []CropType{CropTypeTokens,
        CropTypeInputTextTokens, CropTypeInputAudioTokens,
        CropTypeCachedTextTokens, CropTypeCachedAudioTokens,
        CropTypeOutputTextTokens, CropTypeOutputAudioTokens,
    }

    // UsageTypes*** below exclusive to each other

    UsageTypesDuration = []CropType{
        CropTypeAudioDuration, CropTypeLifetimeAudioDuration,
        CropTypeAudioDurationSingleFeature, CropTypeAudioDurationAllFeatures,
    }
    UsageTypesRequest = []CropType{
        CropTypeRequests,
        CropTypeLifetimeRequestCount,
    }
    UsageTypesSession = []CropType{
        CropTypeSessions,
        CropTypeLifetimeSessionCount,
    }
    UsageTypeTextBytes = []CropType{CropTypeTextBytes}
    UsageTypeTextWords = []CropType{
        CropTypeLifetimeTextWords, CropTypeTextWords,
    }
    UsageTypeTokens    = []CropType{CropTypeLifetimeTokens}
    UsageTypeAllTokens = []CropType{CropTypeTokens,
        CropTypeInputTextTokens, CropTypeInputAudioTokens,
        CropTypeCachedTextTokens, CropTypeCachedAudioTokens,
        CropTypeOutputTextTokens, CropTypeOutputAudioTokens,
    }
)

func (s *CropType) Scan(src any) error {
    if src == nil {
        return nil
    }

    switch value := src.(type) {
    case []byte:
        *s = CropType(value)
    case string:
        *s = CropType(value)
    default:
        return fmt.Errorf("invalid CropType.Scan source type: %T", src)
    }

    return nil
}

func (s CropType) Value() (driver.Value, error) { return string(s), nil }

type CropTypeMeta struct {
    // UnitMeta will be used to calculate the real value of CropType
    UnitMeta CropTypeUnitMeta `json:"unit_meta"`
    // maybe other meta support in the future...
}

func (ct CropTypeMeta) Value() (driver.Value, error) {
    b, err := json.Marshal(ct)
    if err != nil {
        return nil, fmt.Errorf("marshal CropTypeMeta : %w", err)
    }

    return string(b), nil
}

func (ct *CropTypeMeta) Scan(src any) error {
    if src == nil {
        return nil
    }

    srcBytes, ok := src.([]byte)
    if !ok {
        return fmt.Errorf("invalid CropTypeMeta.Scan type: %T", src)
    }

    if len(srcBytes) == 0 {
        return nil
    }

    err := json.Unmarshal(srcBytes, ct)
    if err != nil {
        return fmt.Errorf("unmarshal CropTypeMeta: %w", err)
    }

    return nil
}

// CropTypeUnitMeta only one of DurationUnit & TimesUnit should be set
type CropTypeUnitMeta struct {
    DurationUnit DurationUnit `json:"duration_unit,omitempty"`
    TimesUnit    TimesUnit    `json:"times_unit,omitempty"`
}

const (
    DurationEmpty        = ""
    DurationDefaultAKAMs = ""
    DurationMillisecond  = "millisecond"
    DurationSecond       = "second"
    DurationMinute       = "minute"
    DurationHour         = "hour"

    TimesDefaultAKAOne TimesUnit = ""
    TimesOne           TimesUnit = "one"
    TimesThousand      TimesUnit = "thousand"
    TimesMillion       TimesUnit = "million"
    TimesBillion       TimesUnit = "billion"
)

type DurationUnit string

type TimesUnit string

/*pending_quotas
{
    "service_level": "business",
    "workflow_id": "2192142",
    "quota": [
        {
            "type": "qps",
            "used": 0,
            "limit": 100
        },
        {
            "type": "audio_duration",
            "used": 0,
            "limit": 36000000000
        }
    ]
}
*/

type PendingHarvestLimitations struct {
    Rank          Rank               `json:"service_level"`
    InquiryID     string             `json:"workflow_id"`
    HarvestLimits HarvestLimitations `json:"quota"`
}

func (phl *PendingHarvestLimitations) Scan(src any) error {
    if src == nil {
        return nil
    }

    srcBytes, ok := src.([]byte)
    if !ok {
        return fmt.Errorf("invalid PendingHarvestLimitations.Scan type: %T", src)
    }
    if len(srcBytes) == 0 {
        return nil
    }

    if err := json.Unmarshal(srcBytes, &phl); err != nil {
        return fmt.Errorf("unmarshal PandingHarvestLimitations: %w", err)
    }
    return nil
}

func (phl *PendingHarvestLimitations) Value() (driver.Value, error) {
    b, err := json.Marshal(phl)
    if err != nil {
        return nil, fmt.Errorf("marshal PendingHarvestLimitations: %w", err)
    }

    return b, nil
}
