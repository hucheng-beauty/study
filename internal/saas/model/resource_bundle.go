package model

// BundleResourcePackage is a type of package bundled to other resource package
// so that they have to be ordered at the same time; its fields
// decide when and which the bundle should apply
type BundleResourcePackage struct {
    // NotBundledWhen decides when the bundle package should not be applied
    NotBundledWhen []BundledTiming `json:"not_bundled_when"`

    // RefPreferCode is the prefer code of the package which this is the package referring to.
    RefPreferCode string `json:"ref_prefer_code"`
}

type BundledTiming string
