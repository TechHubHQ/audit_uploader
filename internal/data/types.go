package data

import "time"

// AuditDate formats as dd-mm-yyyy in JSON/log output instead of Go's default RFC3339.
type AuditDate time.Time

func (d AuditDate) MarshalJSON() ([]byte, error) {
	return []byte(`"` + time.Time(d).Format("02-01-2006") + `"`), nil
}

type AuditRecord struct {
	UHID                     string
	DoctorName               string
	Department               string
	AuditDate                AuditDate
	Location                 string
	DrugAllergiesDocumented  *bool
	TotalDrugsInPrescription int
	Prescription             PrescriptionAudit
	Transcription            TranscriptionAudit
	Dispensing               DispensingAudit
	Administration           AdministrationAudit
	AuditObservations        string
}

type AuditCheck struct {
	Answer        *bool
	AffectedDrugs *int
}

type PrescriptionAudit struct {
	Doses                        AuditCheck
	Frequencies                  AuditCheck
	Routes                       AuditCheck
	Units                        AuditCheck
	Concentrations               AuditCheck
	AdministrationRates          AuditCheck
	AllDrugSelectionsAppropriate AuditCheck
	Legibility                   AuditCheck
	StandardAbbreviations        AuditCheck
	CapitalizedDrugNames         AuditCheck
	DrugDrugInteractions         AuditCheck
	DrugFoodInteractions         AuditCheck
}

type TranscriptionAudit struct {
	WrongFormulations AuditCheck
	WrongDrugs        AuditCheck
	WrongStrengths    AuditCheck
}

type DispensingAudit struct {
	CorrectDrugs                           AuditCheck
	CorrectDrugDoses                       AuditCheck
	CorrectDrugFormulations                AuditCheck
	DrugsDispensedBeforeExpiryOrNearExpiry AuditCheck
	CorrectLabelling                       AuditCheck
	WithinDefinedTime                      AuditCheck
	GenericSubstitutionWithoutConsultation AuditCheck
}

type AdministrationAudit struct {
	CorrectPatient                        AuditCheck
	AllDrugsAdministered                  AuditCheck
	CorrectDose                           AuditCheck
	CorrectDrugs                          AuditCheck
	CorrectDosageForm                     AuditCheck
	CorrectRoute                          AuditCheck
	CorrectRate                           AuditCheck
	CorrectDuration                       AuditCheck
	CorrectTime                           AuditCheck
	DrugAdministrationDocumentedProperly  AuditCheck
	NursingDocumentationCompleteAndProper AuditCheck
	DocumentationWithoutAdministration    AuditCheck
}

type RawAuditRecord struct {
	UHID                                                      string
	DoctorName                                                string
	Department                                                string
	AuditDate                                                 string
	Location                                                  string
	DrugAllergiesDocumented                                   string
	TotalDrugsInPrescription                                  string
	AllDrugDosesStatedAppropriately                           string
	DrugsWithInappropriatelyStatedDoses                       string
	AllDrugFrequenciesStatedAppropriately                     string
	DrugsWithInappropriatelyStatedFrequencies                 string
	AllDrugRoutesStatedAppropriately                          string
	DrugsWithInappropriatelyStatedRoutes                      string
	AllDrugUnitsMentionedAppropriately                        string
	DrugsWithUnitsNotMentioned                                string
	AllDrugConcentrationsMentioned                            string
	DrugsWithConcentrationNotMentioned                        string
	AllDrugAdministrationRatesMentioned                       string
	DrugsWithAdministrationRateNotMentioned                   string
	AllDrugSelectionsAppropriate                              string
	DrugsSelectedInappropriately                              string
	PrescriptionLegible                                       string
	IllegibleDrugs                                            string
	OnlyStandardAbbreviationsUsed                             string
	NonApprovedAbbreviationsUsed                              string
	PrescriptionWrittenInCapitalLetters                       string
	DrugNamesNotWrittenInCapitalLetters                       string
	DrugDrugInteractionsMentionedAppropriately                string
	DrugsWithoutDoseModificationForDrugDrugInteractions       string
	DrugFoodInteractionsMentionedAppropriately                string
	DrugsWithoutTimeOrDoseModificationForDrugFoodInteractions string
	WrongFormulationTranscribedOrIndented                     string
	WrongFormulationsTranscribedOrIndented                    string
	WrongDrugTranscribedOrIndented                            string
	WrongDrugsTranscribedOrIndented                           string
	WrongStrengthTranscribedOrIndented                        string
	WrongStrengthsTranscribedOrIndented                       string
	AllDrugsDispensedCorrectly                                string
	WrongDrugsDispensed                                       string
	AllDrugDosesDispensedCorrectly                            string
	WrongDrugDosesDispensed                                   string
	AllDrugFormulationsDispensedCorrectly                     string
	WrongDrugFormulationsDispensed                            string
	DrugsDispensedBeforeExpiry                                string
	ExpiredDrugsDispensed                                     string
	DrugsDispensedWithCorrectLabelling                        string
	DrugsDispensedWithWrongOrNoLabelling                      string
	AllDrugsDispensedWithinDefinedTime                        string
	DrugsNotDispensedWithinDefinedTime                        string
	GenericSubstituteDoneWithoutConsultation                  string
	GenericSubstitutesDoneWithoutConsultation                 string
	DrugsAdministeredToCorrectPatient                         string
	DrugsAdministeredToWrongPatient                           string
	AllDrugsAdministeredToPatient                             string
	DrugsOmittedToPatient                                     string
	AllDrugsAdministeredInCorrectDose                         string
	DrugDosesAdministeredImproperly                           string
	AllDrugsAdministeredCorrectly                             string
	WrongDrugsAdministered                                    string
	AllDrugsAdministeredInCorrectDosageForm                   string
	WrongDosageFormsAdministered                              string
	AllDrugsAdministeredInRightRoute                          string
	DrugsAdministeredInWrongRoute                             string
	AllDrugsAdministeredAtCorrectRate                         string
	DrugsAdministeredAtWrongRate                              string
	AllDrugsAdministeredForCorrectDuration                    string
	DrugsAdministeredForWrongDuration                         string
	AllDrugsAdministeredAtCorrectTime                         string
	DrugsAdministeredAtWrongTime                              string
	DrugAdministrationDocumentedProperly                      string
	DrugsWithoutAdministrationDocumentation                   string
	DrugDocumentationCompleteAndProperByNursingStaff          string
	DrugsDocumentedIncompletelyOrImproperly                   string
	DocumentationWithoutAdministration                        string
	DrugsDocumentedWithoutAdministration                      string
	AuditObservations                                         string
}
