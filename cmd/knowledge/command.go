package main

import (
	"fmt"
	"strings"
	"time"
)

type command struct {
	operation string
	options   []option
}

type option struct {
	name  string
	value string
}

type cliError struct {
	code    errorCode
	message string
	field   string
}

func (e cliError) Error() string {
	return e.message
}

type optionCardinality int

const (
	singleOption optionCardinality = iota
	repeatableOption
)

type operationSpec struct {
	options map[string]optionCardinality
}

var operationSpecs = map[string]operationSpec{
	"search-text": {
		options: map[string]optionCardinality{
			"query": singleOption,
		},
	},
	"search-concept": {
		options: map[string]optionCardinality{
			"concept": singleOption,
		},
	},
	"search-related": {
		options: map[string]optionCardinality{
			"seed-kind":     singleOption,
			"seed-id":       singleOption,
			"relation-type": repeatableOption,
		},
	},
	"get": {
		options: map[string]optionCardinality{
			"assertion-id": singleOption,
		},
	},
	"get-evidence": {
		options: map[string]optionCardinality{
			"assertion-id": singleOption,
		},
	},
	"search-contradictions": {
		options: map[string]optionCardinality{
			"assertion-id": singleOption,
			"concept":      singleOption,
		},
	},
	"search-temporal": {
		options: map[string]optionCardinality{
			"concept":     singleOption,
			"scope-key":   repeatableOption,
			"scope-value": repeatableOption,
			"at":          singleOption,
			"valid-from":  singleOption,
			"valid-until": singleOption,
		},
	},
	"create": {
		options: createOptions(),
	},
	"attach-evidence": {
		options: map[string]optionCardinality{
			"assertion-id":         singleOption,
			"evidence-kind":        singleOption,
			"evidence-text":        singleOption,
			"evidence-observed-at": singleOption,
		},
	},
	"revise": {
		options: reviseOptions(),
	},
	"supersede": {
		options: map[string]optionCardinality{
			"superseded-assertion-id":  singleOption,
			"replacement-assertion-id": singleOption,
		},
	},
}

func revisionOptions() map[string]optionCardinality {
	options := map[string]optionCardinality{
		"valid-from":    singleOption,
		"valid-until":   singleOption,
		"version-scope": singleOption,
		"observed-at":   singleOption,
		"last-verified": singleOption,
		"scope-key":     repeatableOption,
		"scope-value":   repeatableOption,
	}

	return options
}

func createOptions() map[string]optionCardinality {
	options := revisionOptions()
	options["normalized-text"] = singleOption
	options["concept"] = repeatableOption
	options["concept-alias"] = repeatableOption
	options["alias-kind"] = repeatableOption
	options["alias-value"] = repeatableOption
	options["relation-type"] = repeatableOption
	options["relation-target-kind"] = repeatableOption
	options["relation-target-id"] = repeatableOption
	options["evidence-kind"] = repeatableOption
	options["evidence-text"] = repeatableOption
	options["evidence-observed-at"] = repeatableOption

	return options
}

func reviseOptions() map[string]optionCardinality {
	options := revisionOptions()
	options["assertion-id"] = singleOption
	options["normalized-text"] = singleOption

	return options
}

func parseCommand(arguments []string) (command, cliError) {
	if len(arguments) == 0 {
		return command{}, validationFailure("operation", "操作を指定してください")
	}
	specification, ok := operationSpecs[arguments[0]]
	if !ok {
		return command{}, validationFailure("operation", "未知の操作です")
	}
	parsed := command{operation: arguments[0]}
	seen := make(map[string]bool)
	for index := 1; index < len(arguments); index += 2 {
		name, value, err := parseOption(arguments, index)
		if err.code != "" {
			return command{}, err
		}
		cardinality, ok := specification.options[name]
		if !ok {
			return command{}, validationFailure(name, "この操作では使えないoptionです")
		}
		if cardinality == singleOption && seen[name] {
			return command{}, validationFailure(name, "単一値optionを重複指定できません")
		}
		seen[name] = true
		parsed.options = append(parsed.options, option{
			name:  name,
			value: value,
		})
	}
	if err := validateCommand(parsed); err.code != "" {
		return command{}, err
	}

	return parsed, cliError{}
}

func parseOption(arguments []string, index int) (string, string, cliError) {
	if !strings.HasPrefix(arguments[index], "--") || len(arguments[index]) == 2 {
		return "", "", validationFailure("option", "option名を指定してください")
	}
	name := strings.TrimPrefix(arguments[index], "--")
	if index+1 == len(arguments) || strings.HasPrefix(arguments[index+1], "--") {
		return "", "", validationFailure(name, "optionの値を指定してください")
	}

	return name, arguments[index+1], cliError{}
}

func validateCommand(parsed command) cliError {
	values := optionValues(parsed.options)
	if err := validateRequiredOptions(parsed.operation, values); err.code != "" {
		return err
	}
	for _, entry := range parsed.options {
		if strings.TrimSpace(entry.value) == "" {
			return validationFailure(entry.name, "空文字列は指定できません")
		}
	}
	if err := validateGroupOrder(parsed.operation, parsed.options); err.code != "" {
		return err
	}
	if err := validateScopeGroups(values); err.code != "" {
		return err
	}
	if err := validateTemporalOptions(values); err.code != "" {
		return err
	}
	switch parsed.operation {
	case "search-related":
		return validateSearchRelated(values)
	case "search-contradictions":
		return validateOneSelector(values)
	case "search-temporal":
		return validateTemporalSelector(values)
	case "create":
		return validateCreate(values)
	case "attach-evidence":
		return validateEvidenceGroups(values, true)
	case "supersede":
		if values["superseded-assertion-id"][0] == values["replacement-assertion-id"][0] {
			return cliError{
				code:    conflictError,
				message: "同じAssertionを置換できません",
				field:   "replacement-assertion-id",
			}
		}
	}

	return cliError{}
}

func validateGroupOrder(operation string, options []option) cliError {
	for index, entry := range options {
		switch entry.name {
		case "scope-key":
			if !nextOptionIs(options, index, "scope-value") {
				return validationFailure("scope", "Scope groupを順序どおりに指定してください")
			}
		case "alias-kind":
			if operation == "create" && !nextOptionIs(options, index, "alias-value") {
				return validationFailure("alias", "Assertion Alias groupを順序どおりに指定してください")
			}
		case "relation-type":
			if operation == "create" && (!nextOptionIs(options, index, "relation-target-kind") || !nextOptionIs(options, index+1, "relation-target-id")) {
				return validationFailure("relation", "Relation groupを順序どおりに指定してください")
			}
		case "evidence-kind":
			if !nextOptionIs(options, index, "evidence-text") || !nextOptionIs(options, index+1, "evidence-observed-at") {
				return validationFailure("evidence", "Evidence groupを順序どおりに指定してください")
			}
		case "concept-alias":
			if operation == "create" && !hasPreviousConcept(options, index) {
				return validationFailure("concept-alias", "ConceptなしにAliasを指定できません")
			}
		}
	}

	return cliError{}
}

func nextOptionIs(options []option, index int, name string) bool {
	return index+1 < len(options) && options[index+1].name == name
}

func hasPreviousConcept(options []option, index int) bool {
	for position := index - 1; position >= 0; position-- {
		if options[position].name == "concept" {
			return true
		}
	}

	return false
}

func optionValues(options []option) map[string][]string {
	values := make(map[string][]string)
	for _, entry := range options {
		values[entry.name] = append(values[entry.name], entry.value)
	}

	return values
}

func validateRequiredOptions(operation string, values map[string][]string) cliError {
	required := map[string][]string{
		"search-text":     {"query"},
		"search-concept":  {"concept"},
		"search-related":  {"seed-kind", "seed-id"},
		"get":             {"assertion-id"},
		"get-evidence":    {"assertion-id"},
		"attach-evidence": {"assertion-id", "evidence-kind", "evidence-text", "evidence-observed-at"},
		"revise":          {"assertion-id", "normalized-text"},
		"supersede":       {"superseded-assertion-id", "replacement-assertion-id"},
		"create":          {"normalized-text"},
	}
	for _, name := range required[operation] {
		if len(values[name]) == 0 {
			return validationFailure(name, "必須optionを指定してください")
		}
	}

	return cliError{}
}

func validateScopeGroups(values map[string][]string) cliError {
	keys := values["scope-key"]
	entries := values["scope-value"]
	if len(keys) != len(entries) {
		return validationFailure("scope", "Scope groupを完全に指定してください")
	}
	seen := make(map[string]bool)
	for _, key := range keys {
		if seen[key] {
			return validationFailure("scope-key", "Scope keyを重複指定できません")
		}
		seen[key] = true
	}

	return cliError{}
}

func validateTemporalOptions(values map[string][]string) cliError {
	for _, name := range []string{
		"at",
		"valid-from",
		"valid-until",
		"observed-at",
		"last-verified",
		"evidence-observed-at",
	} {
		for _, value := range values[name] {
			parsed, err := time.Parse(time.RFC3339, value)
			if err != nil {
				return validationFailure(name, "RFC 3339 UTC時刻を指定してください")
			}
			_, offset := parsed.Zone()
			if offset != 0 {
				return validationFailure(name, "UTC時刻を指定してください")
			}
		}
	}
	from := values["valid-from"]
	until := values["valid-until"]
	if len(from) == 1 && len(until) == 1 {
		fromTime, _ := time.Parse(time.RFC3339, from[0])
		untilTime, _ := time.Parse(time.RFC3339, until[0])
		if fromTime.After(untilTime) {
			return validationFailure("valid-until", "valid-from以前の時刻を指定してください")
		}
	}

	return cliError{}
}

func validateSearchRelated(values map[string][]string) cliError {
	if !oneOf(values["seed-kind"][0], "assertion", "concept") {
		return validationFailure("seed-kind", "許可されていない種別です")
	}
	for _, relationType := range values["relation-type"] {
		if !oneOf(relationType, "related_to", "prerequisite", "causes", "contributes_to", "contradicts", "supersedes") {
			return validationFailure("relation-type", "許可されていないRelation種別です")
		}
	}

	return cliError{}
}

func validateOneSelector(values map[string][]string) cliError {
	if len(values["assertion-id"])+len(values["concept"]) != 1 {
		return validationFailure("selector", "AssertionまたはConceptを一つだけ指定してください")
	}

	return cliError{}
}

func validateTemporalSelector(values map[string][]string) cliError {
	if len(values["concept"]) == 0 && len(values["scope-key"]) == 0 {
		return validationFailure("selector", "ConceptまたはScopeを指定してください")
	}
	if len(values["at"]) > 0 && (len(values["valid-from"]) > 0 || len(values["valid-until"]) > 0) {
		return validationFailure("at", "atとvalid-from、valid-untilは同時に指定できません")
	}
	if len(values["at"]) == 0 && len(values["valid-from"]) != len(values["valid-until"]) {
		return validationFailure("valid-from", "valid-fromとvalid-untilは組で指定してください")
	}

	return cliError{}
}

func validateCreate(values map[string][]string) cliError {
	if err := validateConceptGroups(values); err.code != "" {
		return err
	}
	if err := validateAliasGroups(values); err.code != "" {
		return err
	}
	if err := validateRelationGroups(values); err.code != "" {
		return err
	}

	return validateEvidenceGroups(values, false)
}

func validateConceptGroups(values map[string][]string) cliError {
	if err := validateUnique(values["concept"], "concept", "Conceptを重複指定できません"); err.code != "" {
		return err
	}

	return validateUnique(values["concept-alias"], "concept-alias", "Concept Aliasを重複指定できません")
}

func validateAliasGroups(values map[string][]string) cliError {
	kinds := values["alias-kind"]
	entries := values["alias-value"]
	seen := make(map[string]bool)
	for index, kind := range kinds {
		if !oneOf(kind, "api_name", "identifier") {
			return validationFailure("alias-kind", "許可されていないAlias種別です")
		}
		key := kind + "\x00" + entries[index]
		if seen[key] {
			return validationFailure("alias", "Assertion Aliasを重複指定できません")
		}
		seen[key] = true
	}

	return cliError{}
}

func validateRelationGroups(values map[string][]string) cliError {
	types := values["relation-type"]
	kinds := values["relation-target-kind"]
	for index, relationType := range types {
		if !oneOf(relationType, "related_to", "prerequisite", "causes", "contributes_to", "contradicts") {
			return validationFailure("relation-type", "許可されていないRelation種別です")
		}
		if !oneOf(kinds[index], "assertion", "concept") {
			return validationFailure("relation-target-kind", "許可されていないRelation対象種別です")
		}
		if relationType == "contradicts" && kinds[index] != "assertion" {
			return validationFailure("relation-target-kind", "contradictsの対象はAssertionです")
		}
	}

	return cliError{}
}

func validateEvidenceGroups(values map[string][]string, exactlyOne bool) cliError {
	kinds := values["evidence-kind"]
	texts := values["evidence-text"]
	observed := values["evidence-observed-at"]
	if len(kinds) != len(texts) || len(kinds) != len(observed) || (exactlyOne && len(kinds) != 1) || (!exactlyOne && len(kinds) == 0) {
		return validationFailure("evidence", "Evidence groupを完全に指定してください")
	}
	seen := make(map[string]bool)
	for index, kind := range kinds {
		if !oneOf(kind, "user_explanation", "user_reasoning", "user_code", "self_report", "concept_recognition", "correction") {
			return validationFailure("evidence-kind", "許可されていないEvidence種別です")
		}
		key := kind + "\x00" + texts[index] + "\x00" + observed[index]
		if seen[key] {
			return validationFailure("evidence", "Evidenceを重複指定できません")
		}
		seen[key] = true
	}

	return cliError{}
}

func validateUnique(values []string, field string, message string) cliError {
	seen := make(map[string]bool)
	for _, value := range values {
		if seen[value] {
			return validationFailure(field, message)
		}
		seen[value] = true
	}

	return cliError{}
}

func oneOf(value string, allowed ...string) bool {
	for _, candidate := range allowed {
		if value == candidate {
			return true
		}
	}

	return false
}

func validationFailure(field string, message string) cliError {
	return cliError{
		code:    validationError,
		message: fmt.Sprintf("%s: %s", field, message),
		field:   field,
	}
}
