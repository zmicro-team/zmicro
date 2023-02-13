package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"

	"github.com/zmicro-team/zmicro/core/enumerate"
)

func runProtoGen(gen *protogen.Plugin) error {
	var mergeEnums []*Enum
	var source []string
	gen.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)

	isMerge := *merge
	if *merge {
		if *_package == "" ||
			*filename == "" ||
			*goPackage == "" {
			return errors.New("when enable merge, filename,package,go_package must be set")
		}
		mergeEnums = make([]*Enum, 0, len(gen.Files)*4)
		source = make([]string, 0, len(gen.Files))
	}
	usedTemplate := enumTemplate
	if *customTemplate != "" {
		t, err := ParseTemplateFromFile(*customTemplate)
		if err != nil {
			return err
		}
		usedTemplate = t
	}

	for _, f := range gen.Files {
		if !f.Generate {
			continue
		}
		enums := intoEnums("", f.Enums)
		enums = append(enums, intoEnumsFromMessage("", f.Messages)...)
		if len(enums) == 0 {
			continue
		}
		if isMerge {
			source = append(source, f.Desc.Path())
			mergeEnums = append(mergeEnums, enums...)
			continue
		}
		g := gen.NewGeneratedFile(f.GeneratedFilenamePrefix+*suffix, f.GoImportPath)
		e := &EnumFile{
			Version:       version,
			ProtocVersion: protocVersion(gen),
			IsDeprecated:  f.Proto.GetOptions().GetDeprecated(),
			Source:        f.Desc.Path(),
			Package:       string(f.GoPackageName),
			Enums:         enums,
		}
		_ = e.execute(usedTemplate, g)
	}
	if isMerge {
		mergeFile := &EnumFile{
			Version:       version,
			ProtocVersion: protocVersion(gen),
			IsDeprecated:  false,
			Source:        strings.Join(source, ","),
			Package:       *_package,
			Enums:         mergeEnums,
		}
		g := gen.NewGeneratedFile(*filename+*suffix, protogen.GoImportPath(*goPackage))
		return mergeFile.execute(usedTemplate, g)
	}
	return nil
}

func protocVersion(gen *protogen.Plugin) string {
	v := gen.Request.GetCompilerVersion()
	if v == nil {
		return "(unknown)"
	}
	var suffix string
	if s := v.GetSuffix(); s != "" {
		suffix = "-" + s
	}
	return fmt.Sprintf("v%d.%d.%d%s", v.GetMajor(), v.GetMinor(), v.GetPatch(), suffix)
}

// intoEnumsFromMessage generates the errors definitions, excluding the package statement.
func intoEnumsFromMessage(nestedMessageName string, protoMessages []*protogen.Message) []*Enum {
	enums := make([]*Enum, 0, 128)
	for _, pm := range protoMessages {
		tmpNestedMessageName := string(pm.Desc.Name())
		if nestedMessageName != "" {
			tmpNestedMessageName = nestedMessageName + "_" + tmpNestedMessageName
		}
		enums = append(enums, intoEnums(tmpNestedMessageName, pm.Enums)...)
		enums = append(enums, intoEnumsFromMessage(tmpNestedMessageName, pm.Messages)...)
	}
	return enums
}

// intoEnums generates the errors definitions, excluding the package statement.
func intoEnums(nestedMessageName string, protoEnums []*protogen.Enum) []*Enum {
	enums := make([]*Enum, 0, len(protoEnums))
	for _, pe := range protoEnums {
		if len(pe.Values) == 0 {
			continue
		}
		isEnabled := proto.GetExtension(pe.Desc.Options(), enumerate.E_Enabled)
		ok := isEnabled.(bool)
		if !ok {
			continue
		}

		eValueMp := make(map[int]string, len(pe.Values))
		eValues := make([]*EnumValue, 0, len(pe.Values))
		for _, v := range pe.Values {
			mpv := proto.GetExtension(v.Desc.Options(), enumerate.E_Mapping)
			mappingValue, _ := mpv.(string)

			eValues = append(eValues, &EnumValue{
				Value:      string(v.Desc.Name()),
				Number:     int(v.Desc.Number()),
				CamelValue: CamelCase(string(v.Desc.Name())),
				Mapping:    mappingValue,
				Comment:    strings.TrimSuffix(string(v.Comments.Leading), "\n"),
			})
			eValueMp[v.Desc.Index()] = mappingValue
		}
		b, _ := json.Marshal(eValueMp)
		bb := strings.ReplaceAll(string(b), `"`, "")
		bb = strings.Replace(bb, "{", "[", 1)
		bb = strings.Replace(bb, "}", "]", 1)
		name := string(pe.Desc.Name())
		if nestedMessageName != "" {
			name = nestedMessageName + "_" + name
		}
		enums = append(enums, &Enum{
			Name:    name,
			Comment: strings.ReplaceAll(string(pe.Comments.Leading), "\n", "") + ", " + bb,
			Values:  eValues,
		})
	}
	return enums
}
