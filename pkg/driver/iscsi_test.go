package driver

import "testing"

// buildConnector must set the fields csi-lib-iscsi actually consults to configure
// CHAP: Secrets.SecretsType == "chap" and Connector.DoCHAPDiscovery. Without them
// the library writes no CHAP settings and a CHAP-protected portal rejects login.
func TestBuildConnector_CHAPEnabled(t *testing.T) {
	h := &ISCSIHandler{}
	config := &ISCSIConfig{
		TargetPortal:   "10.0.0.1:3260",
		TargetIQN:      "iqn.2000-01.io.truenas:vol",
		CHAPUsername:   "k8s",
		CHAPPassword:   "pass1234abcd",
		CHAPUsernameIn: "truenas",
		CHAPPasswordIn: "peerpass1234",
	}

	c := h.buildConnector("vol1", config)

	if !c.DoCHAPDiscovery {
		t.Error("DoCHAPDiscovery must be true when CHAP is configured")
	}
	if c.DiscoverySecrets.SecretsType != "chap" || c.SessionSecrets.SecretsType != "chap" {
		t.Errorf("SecretsType must be \"chap\": discovery=%q session=%q",
			c.DiscoverySecrets.SecretsType, c.SessionSecrets.SecretsType)
	}
	if c.SessionSecrets.UserName != "k8s" || c.SessionSecrets.Password != "pass1234abcd" {
		t.Errorf("session credentials not propagated: %+v", c.SessionSecrets)
	}
	if c.SessionSecrets.UserNameIn != "truenas" || c.SessionSecrets.PasswordIn != "peerpass1234" {
		t.Errorf("mutual (incoming) credentials not propagated: %+v", c.SessionSecrets)
	}
}

func TestBuildConnector_NoCHAP(t *testing.T) {
	h := &ISCSIHandler{}
	c := h.buildConnector("vol1", &ISCSIConfig{
		TargetPortal: "10.0.0.1:3260",
		TargetIQN:    "iqn.2000-01.io.truenas:vol",
	})

	if c.DoCHAPDiscovery {
		t.Error("DoCHAPDiscovery must be false when no CHAP credentials are set")
	}
	if c.DiscoverySecrets.SecretsType != "" {
		t.Errorf("SecretsType should be empty without CHAP, got %q", c.DiscoverySecrets.SecretsType)
	}
}
