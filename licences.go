package irapi

type LicenceClass int

const (
	LicenceClassRookie LicenceClass = iota + 1
	LicenceClassD
	LicenceClassC
	LicenceClassB
	LicenceClassA
	LicenceClassPro
	LicenceClassProWC

	licenseClassStrings = " RDCBDA"
)

func (l LicenceClass) String() string {
	if l == LicenceClassProWC {
		return "Pro/WC"
	}

	if l == LicenceClassPro {
		return "Pro"
	}

	return string(licenseClassStrings[l])
}

type LicenceCategory int

const (
	LicenceCategoryRoad LicenceCategory = iota + 1
	LicenceCategoryOval
	LicenceCategoryDirtRoad
	LicenceCategoryDirtOval
)

func (l LicenceCategory) String() string {
	s := ""

	if l >= LicenceCategoryDirtRoad {
		s = "Dirt "
	}

	if l%2 == 0 {
		s += "Road"
	} else {
		s += "Oval"
	}

	return s
}
