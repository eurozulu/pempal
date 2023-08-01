package prompts

import (
	"fmt"
	"github.com/eurozulu/pempal/builders"
	"github.com/eurozulu/pempal/resources"
	"github.com/eurozulu/pempal/ui"
	"github.com/eurozulu/pempal/utils"
)

const generageKeyText = "Generate new key"

func supportseResourceTypeView(resourceType resources.ResourceType) error {
	switch resourceType {
	case resources.PrivateKey, resources.CertificateRequest, resources.Certificate, resources.RevocationList:
		return nil
	default:
		return fmt.Errorf("no resource view available for %s", resourceType.String())
	}
}

func createResourceTypeView(resourceType resources.ResourceType) (ui.ParentView, error) {
	switch resourceType {
	case resources.PrivateKey:
		return privateKeyView(), nil
	case resources.CertificateRequest:
		return requestView(), nil
	case resources.Certificate:
		return certificateView(), nil
	case resources.RevocationList:
		return revokationListView(), nil

	default:
		return nil, fmt.Errorf("no resource view available for %s", resourceType.String())
	}
}

func privateKeyView() ui.ParentView {
	return ui.NewParentView("Private Key", "", []ui.View{
		ui.NewTextList(builders.Property_key_algorithm, "RSA", utils.PublicKeyAlgorithms...),
		ui.NewTextListHidden(builders.Property_key_length, "", []string{"512", "1024", "2048", "4096"}...),
		ui.NewTextListHidden(builders.Property_key_curve, "", utils.ECDSACurveNames[1:]...),
		ui.NewBoolViewPreSelected(builders.Property_key_is_encrypted, true),
	}...)
}

func requestView() ui.ParentView {
	return ui.NewParentView("Certificate Request", "", []ui.View{
		ui.NewTextView(builders.Property_signature, ""),
		ui.NewTextList(builders.Property_signature_algorithm, "SHA512-RSA", utils.SignatureAlgorithmNames()...),
		ui.NewTextList(builders.Property_key_algorithm, "RSA", utils.PublicKeyAlgorithms...),
		ui.NewTextList(builders.Property_public_key, "", "", generageKeyText),
		NewDNView(builders.Property_subject, ""),
	}...)
}

func certificateView() ui.ParentView {
	return ui.NewParentView("Certificate", "", []ui.View{
		ui.NewNumberView(builders.Property_version, 0),
		ui.NewNumberView(builders.Property_serial_number, 0),
		ui.NewTextView(builders.Property_signature, ""),
		ui.NewTextList(builders.Property_signature_algorithm, "SHA512-RSA", utils.SignatureAlgorithmNames()...),
		ui.NewTextList(builders.Property_key_algorithm, "RSA", utils.PublicKeyAlgorithms...),
		NewPublicKeyView(builders.Property_public_key, ""),
		NewDNView(builders.Property_issuer, ""),
		NewDNView(builders.Property_subject, ""),
		NewDateView(builders.Property_not_before, ""),
		NewDateView(builders.Property_not_after, ""),
		ui.NewBoolViewPreSelected(builders.Property_is_ca, false),
		ui.NewBoolView(builders.Property_basic_constraints_valid),
		ui.NewNumberView(builders.Property_max_path_len, 0),
		ui.NewBoolView(builders.Property_max_path_len_zero),
		ui.NewMultiSelectHidden(builders.Property_key_usage, "", utils.KeyUsageNames...),
		ui.NewMultiSelectHidden(builders.Property_extended_key_usage, "", utils.ExtKeyUsageNames...),
	}...)
}

func revokationListView() ui.ParentView {
	return ui.NewParentView("Revokation List", "", []ui.View{
		NewDNView(builders.Property_issuer, ""),
		ui.NewTextView(builders.Property_signature, ""),
		ui.NewTextList(builders.Property_signature_algorithm, "SHA512-RSA", utils.SignatureAlgorithmNames()...),

		ui.NewTextList(builders.Property_revoked_certificates, "SHA512-RSA", utils.SignatureAlgorithmNames()...),
		ui.NewNumberView(builders.Property_number, 0),
		NewDateView(builders.Property_this_update, ""),
		NewDateView(builders.Property_next_update, ""),
		ui.NewMultiSelectHidden(builders.Property_extensions, "", utils.KeyUsageNames...),
		ui.NewMultiSelectHidden(builders.Property_extra_extensions, "", utils.KeyUsageNames...),
		ui.NewTextListHidden(builders.Property_revokation_list, ""),
	}...)
}
