package model

// AliasGenerator defines how resource pack alias is generated
// following the blueprint.
type AliasGenerator string

const (
    DefaultAliasGenerator       AliasGenerator = ""
    EmptyAliasGenerator         AliasGenerator = "empty"
    LeafDisplayAliasGenerator   AliasGenerator = "leaf_display"
    PreferCodeAliasGenerator    AliasGenerator = "prefer_code"
    SpecifiedNameAliasGenerator AliasGenerator = "specified_name"

    leafDisplayAliasDelimiter = '-'
    leadDisplayDetailsKey     = "display"
)
