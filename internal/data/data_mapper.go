package data

import (
	"audituploader/internal/log"
	"strconv"
	"strings"
	"time"
)

func buildRawAuditRecord(rows map[string]string) RawAuditRecord {
	rawRecord := RawAuditRecord{
		UHID:                                  rows["uhid/ipnumber"],
		DoctorName:                            rows["doctorname"],
		Department:                            rows["department"],
		AuditDate:                             rows["auditdate"],
		Location:                              rows["location"],
		DrugAllergiesDocumented:               rows["drugallergiesdocumented"],
		TotalDrugsInPrescription:              rows["totalnumberofdrugsintheprescription"],
		AllDrugDosesStatedAppropriately:       rows["werealldrugsdosesstatedappropriately"],
		DrugsWithInappropriatelyStatedDoses:   rows["howmanydrugsdidnothavedosesstatedappropriately"],
		AllDrugFrequenciesStatedAppropriately: rows["werealldrugsfrequencystatedappropriately"],
		DrugsWithInappropriatelyStatedFrequencies:                 rows["howmanydrugsdidnothavefrequencystatedappropriately"],
		AllDrugRoutesStatedAppropriately:                          rows["werealldrugsroutesstatedappropriately"],
		DrugsWithInappropriatelyStatedRoutes:                      rows["howmanydrugsdidnohaveroutestatedapproriately"],
		AllDrugUnitsMentionedAppropriately:                        rows["werealldrugsarehavingunitsmentionedappropriately"],
		DrugsWithUnitsNotMentioned:                                rows["howmanydrugsnotmentionedunits"],
		AllDrugConcentrationsMentioned:                            rows["werealldrugsarehavingconcentrationmentioned"],
		DrugsWithConcentrationNotMentioned:                        rows["howmanydrugsnotmentionedconcentration"],
		AllDrugAdministrationRatesMentioned:                       rows["werealldrugshavingrateofadministrationmentioned"],
		DrugsWithAdministrationRateNotMentioned:                   rows["howmanydrugsdidnothaverateofconcentrationmentioned"],
		AllDrugSelectionsAppropriate:                              rows["werealldrugsselectionwasappropriate"],
		DrugsSelectedInappropriately:                              rows["howmanydrugswereselectedinappropriately"],
		PrescriptionLegible:                                       rows["wastheprescriptionlegible"],
		IllegibleDrugs:                                            rows["howmanydrugsareillegible"],
		OnlyStandardAbbreviationsUsed:                             rows["wereonlystandardabbreviationsusedintheprescripton"],
		NonApprovedAbbreviationsUsed:                              rows["howmanynonapprovedabbreviationsused"],
		PrescriptionWrittenInCapitalLetters:                       rows["wastheprescriptionwrittenincapitalletters"],
		DrugNamesNotWrittenInCapitalLetters:                       rows["howmanydrugnameswerenotwrittenincapitalletters"],
		DrugDrugInteractionsMentionedAppropriately:                rows["weredrug-druginteractionsmentionedappropriately"],
		DrugsWithoutDoseModificationForDrugDrugInteractions:       rows["howmanydrugshadnonmodificationofdrugdosekeepinginminddrugdruginteractions"],
		DrugFoodInteractionsMentionedAppropriately:                rows["weredrug-foodinteractionsmentionedappropriately"],
		DrugsWithoutTimeOrDoseModificationForDrugFoodInteractions: rows["howmanydrugshadnonmodificationoftimeofadministrationordosekeepinginminddrugfoodinteractions"],
		WrongFormulationTranscribedOrIndented:                     rows["wrongformulationtranscribed/indented"],
		WrongFormulationsTranscribedOrIndented:                    rows["howmanywrongformulationtranscribed/indented"],
		WrongDrugTranscribedOrIndented:                            rows["wrongdrugtranscribed/indented"],
		WrongDrugsTranscribedOrIndented:                           rows["howmanywrongdrugtranscribed/indented"],
		WrongStrengthTranscribedOrIndented:                        rows["wrongstrengthtranscribed/indented"],
		WrongStrengthsTranscribedOrIndented:                       rows["howmanywrongstrengthtranscribed/indented"],
		AllDrugsDispensedCorrectly:                                rows["werealldrugsdispensedcorrectly"],
		WrongDrugsDispensed:                                       rows["howmanywrongdrugsdispensed"],
		AllDrugDosesDispensedCorrectly:                            rows["werealldrugdosesdispensedcorrectly"],
		WrongDrugDosesDispensed:                                   rows["howmanywrongdosesdispensed"],
		AllDrugFormulationsDispensedCorrectly:                     rows["werealldrugformulationsdispensedcorrectly"],
		WrongDrugFormulationsDispensed:                            rows["howmanywrongdrugformulationsdispensed"],
		DrugsDispensedBeforeExpiry:                                rows["weredrugsdispensedbeforeexpiry"],
		ExpiredDrugsDispensed:                                     rows["howmanyexpireddrugsdispensed"],
		DrugsDispensedWithCorrectLabelling:                        rows["weredrugsdispensedwithcorrectlabelling"],
		DrugsDispensedWithWrongOrNoLabelling:                      rows["howmanydrugsdispensedinwrong/nodruglabelling"],
		AllDrugsDispensedWithinDefinedTime:                        rows["werealldrugsdispensedwithindefinedtime"],
		DrugsNotDispensedWithinDefinedTime:                        rows["howmanydrugswerenotdispensedindefinedtime"],
		GenericSubstituteDoneWithoutConsultation:                  rows["genericsubstitutedonewithoutconsultation"],
		GenericSubstitutesDoneWithoutConsultation:                 rows["howmanygenericsubstitutedonewithoutconsultation"],
		DrugsAdministeredToCorrectPatient:                         rows["weredrugsadministeredtothecorrectpatient"],
		DrugsAdministeredToWrongPatient:                           rows["howmanydrugswereadministeredtowrongpatient"],
		AllDrugsAdministeredToPatient:                             rows["werealldrugsadministeredtothepatient"],
		DrugsOmittedToPatient:                                     rows["howmanydrugswereomittedtothepatient"],
		AllDrugsAdministeredInCorrectDose:                         rows["werealldrugsadministeredincorrectdose"],
		DrugDosesAdministeredImproperly:                           rows["howmanydrugdoseswereadministeredimproperly"],
		AllDrugsAdministeredCorrectly:                             rows["werealldrugsadministeredcorrectly"],
		WrongDrugsAdministered:                                    rows["howmanywrongdrugswereadminstered"],
		AllDrugsAdministeredInCorrectDosageForm:                   rows["werealldrugsadministeredincorrectdosageform"],
		WrongDosageFormsAdministered:                              rows["howmanywrongdosageformadministered"],
		AllDrugsAdministeredInRightRoute:                          rows["werealldrugsadministeredinrightroute"],
		DrugsAdministeredInWrongRoute:                             rows["howmanywrongrouteofdrugsadministered"],
		AllDrugsAdministeredAtCorrectRate:                         rows["werealldrugsadministeredincorrectrate"],
		DrugsAdministeredAtWrongRate:                              rows["howmanydrugswereadministeredinwrongrate"],
		AllDrugsAdministeredForCorrectDuration:                    rows["werealldrugsadministeredincorrectduration"],
		DrugsAdministeredForWrongDuration:                         rows["howmanydrugswereadministeredinwrongduration"],
		AllDrugsAdministeredAtCorrectTime:                         rows["werealldrugsadministeredincorrecttime"],
		DrugsAdministeredAtWrongTime:                              rows["howmanydrugswereadministeredinwrongtime"],
		DrugAdministrationDocumentedProperly:                      rows["wasdocumentationofdrugadministrationdoneproperly"],
		DrugsWithoutAdministrationDocumentation:                   rows["howmanydrugswerenotdocumentedtheadministration"],
		DrugDocumentationCompleteAndProperByNursingStaff:          rows["wasdocumentationofdrugscompletely&properlydonebynursingstaff"],
		DrugsDocumentedIncompletelyOrImproperly:                   rows["howmanydrugsweredocumentedincompletly/improperly"],
		DocumentationWithoutAdministration:                        rows["documentationwithoutadministration"],
		DrugsDocumentedWithoutAdministration:                      rows["howmanydrugsweredocumentedwithoutadministration"],
		AuditObservations:                                         rows["auditobservations"],
	}
	log.Debug("Raw Audit Record", "record", rawRecord)
	return rawRecord
}

func MapToAuditRecord(raw RawAuditRecord) AuditRecord {
	return AuditRecord{
		UHID:                     raw.UHID,
		DoctorName:               raw.DoctorName,
		Department:               raw.Department,
		AuditDate:                parseAuditDate(raw.AuditDate),
		Location:                 raw.Location,
		DrugAllergiesDocumented:  parseBoolPtr(raw.DrugAllergiesDocumented),
		TotalDrugsInPrescription: parseInt(raw.TotalDrugsInPrescription),
		Prescription:             mapPrescriptionAudit(raw),
		Transcription:            mapTranscriptionAudit(raw),
		Dispensing:               mapDispensingAudit(raw),
		Administration:           mapAdministrationAudit(raw),
		AuditObservations:        raw.AuditObservations,
	}
}

func mapPrescriptionAudit(raw RawAuditRecord) PrescriptionAudit {
	return PrescriptionAudit{
		Doses:                        mapAuditCheck(raw.AllDrugDosesStatedAppropriately, raw.DrugsWithInappropriatelyStatedDoses),
		Frequencies:                  mapAuditCheck(raw.AllDrugFrequenciesStatedAppropriately, raw.DrugsWithInappropriatelyStatedFrequencies),
		Routes:                       mapAuditCheck(raw.AllDrugRoutesStatedAppropriately, raw.DrugsWithInappropriatelyStatedRoutes),
		Units:                        mapAuditCheck(raw.AllDrugUnitsMentionedAppropriately, raw.DrugsWithUnitsNotMentioned),
		Concentrations:               mapAuditCheck(raw.AllDrugConcentrationsMentioned, raw.DrugsWithConcentrationNotMentioned),
		AdministrationRates:          mapAuditCheck(raw.AllDrugAdministrationRatesMentioned, raw.DrugsWithAdministrationRateNotMentioned),
		AllDrugSelectionsAppropriate: mapAuditCheck(raw.AllDrugSelectionsAppropriate, raw.DrugsSelectedInappropriately),
		Legibility:                   mapAuditCheck(raw.PrescriptionLegible, raw.IllegibleDrugs),
		StandardAbbreviations:        mapAuditCheck(raw.OnlyStandardAbbreviationsUsed, raw.NonApprovedAbbreviationsUsed),
		CapitalizedDrugNames:         mapAuditCheck(raw.PrescriptionWrittenInCapitalLetters, raw.DrugNamesNotWrittenInCapitalLetters),
		DrugDrugInteractions:         mapAuditCheck(raw.DrugDrugInteractionsMentionedAppropriately, raw.DrugsWithoutDoseModificationForDrugDrugInteractions),
		DrugFoodInteractions:         mapAuditCheck(raw.DrugFoodInteractionsMentionedAppropriately, raw.DrugsWithoutTimeOrDoseModificationForDrugFoodInteractions),
	}
}

func mapTranscriptionAudit(raw RawAuditRecord) TranscriptionAudit {
	return TranscriptionAudit{
		WrongFormulations: mapAuditCheck(raw.WrongFormulationTranscribedOrIndented, raw.WrongFormulationsTranscribedOrIndented),
		WrongDrugs:        mapAuditCheck(raw.WrongDrugTranscribedOrIndented, raw.WrongDrugsTranscribedOrIndented),
		WrongStrengths:    mapAuditCheck(raw.WrongStrengthTranscribedOrIndented, raw.WrongStrengthsTranscribedOrIndented),
	}
}

func mapDispensingAudit(raw RawAuditRecord) DispensingAudit {
	return DispensingAudit{
		CorrectDrugs:                           mapAuditCheck(raw.AllDrugsDispensedCorrectly, raw.WrongDrugsDispensed),
		CorrectDrugDoses:                       mapAuditCheck(raw.AllDrugDosesDispensedCorrectly, raw.WrongDrugDosesDispensed),
		CorrectDrugFormulations:                mapAuditCheck(raw.AllDrugFormulationsDispensedCorrectly, raw.WrongDrugFormulationsDispensed),
		DrugsDispensedBeforeExpiryOrNearExpiry: mapAuditCheck(raw.DrugsDispensedBeforeExpiry, raw.ExpiredDrugsDispensed),
		CorrectLabelling:                       mapAuditCheck(raw.DrugsDispensedWithCorrectLabelling, raw.DrugsDispensedWithWrongOrNoLabelling),
		WithinDefinedTime:                      mapAuditCheck(raw.AllDrugsDispensedWithinDefinedTime, raw.DrugsNotDispensedWithinDefinedTime),
		GenericSubstitutionWithoutConsultation: mapAuditCheck(raw.GenericSubstituteDoneWithoutConsultation, raw.GenericSubstitutesDoneWithoutConsultation),
	}
}

func mapAdministrationAudit(raw RawAuditRecord) AdministrationAudit {
	return AdministrationAudit{
		CorrectPatient:                        mapAuditCheck(raw.DrugsAdministeredToCorrectPatient, raw.DrugsAdministeredToWrongPatient),
		AllDrugsAdministered:                  mapAuditCheck(raw.AllDrugsAdministeredToPatient, raw.DrugsOmittedToPatient),
		CorrectDose:                           mapAuditCheck(raw.AllDrugsAdministeredInCorrectDose, raw.DrugDosesAdministeredImproperly),
		CorrectDrugs:                          mapAuditCheck(raw.AllDrugsAdministeredCorrectly, raw.WrongDrugsAdministered),
		CorrectDosageForm:                     mapAuditCheck(raw.AllDrugsAdministeredInCorrectDosageForm, raw.WrongDosageFormsAdministered),
		CorrectRoute:                          mapAuditCheck(raw.AllDrugsAdministeredInRightRoute, raw.DrugsAdministeredInWrongRoute),
		CorrectRate:                           mapAuditCheck(raw.AllDrugsAdministeredAtCorrectRate, raw.DrugsAdministeredAtWrongRate),
		CorrectDuration:                       mapAuditCheck(raw.AllDrugsAdministeredForCorrectDuration, raw.DrugsAdministeredForWrongDuration),
		CorrectTime:                           mapAuditCheck(raw.AllDrugsAdministeredAtCorrectTime, raw.DrugsAdministeredAtWrongTime),
		DrugAdministrationDocumentedProperly:  mapAuditCheck(raw.DrugAdministrationDocumentedProperly, raw.DrugsWithoutAdministrationDocumentation),
		NursingDocumentationCompleteAndProper: mapAuditCheck(raw.DrugDocumentationCompleteAndProperByNursingStaff, raw.DrugsDocumentedIncompletelyOrImproperly),
		DocumentationWithoutAdministration:    mapAuditCheck(raw.DocumentationWithoutAdministration, raw.DrugsDocumentedWithoutAdministration),
	}
}

func mapAuditCheck(answerRaw, affectedRaw string) AuditCheck {
	return AuditCheck{
		Answer:        parseBoolPtr(answerRaw),
		AffectedDrugs: parseIntPtr(affectedRaw),
	}
}

func parseAuditDate(raw string) AuditDate {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return AuditDate(time.Time{})
	}

	layouts := []string{"2-Jan-06", "02-Jan-06", "01-02-06", "02-01-06", "2006-01-02", "02/01/2006", "01/02/2006"}
	for _, layout := range layouts {
		if parsed, err := time.Parse(layout, raw); err == nil {
			return AuditDate(parsed)
		}
	}

	log.Error("Unable to parse audit date", "value", raw)
	return AuditDate(time.Time{})
}

func parseBoolPtr(raw string) *bool {
	value := strings.ToLower(strings.TrimSpace(raw))
	switch value {
	case "yes", "y", "true", "1":
		v := true
		return &v
	case "no", "n", "false", "0":
		v := false
		return &v
	default:
		return nil
	}
}

func parseInt(raw string) int {
	if v := parseIntPtr(raw); v != nil {
		return *v
	}
	return 0
}

func parseIntPtr(raw string) *int {
	value := strings.TrimSpace(raw)
	if value == "" {
		return nil
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return nil
	}

	return &parsed
}
