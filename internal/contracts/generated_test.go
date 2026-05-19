package contracts_test

import (
	"testing"

	builderv1 "gitlab.com/Geneden/phanes/gen/go/phanes/builder/v1"
	cachev1 "gitlab.com/Geneden/phanes/gen/go/phanes/cache/v1"
	commonv1 "gitlab.com/Geneden/phanes/gen/go/phanes/common/v1"
	launcherv1 "gitlab.com/Geneden/phanes/gen/go/phanes/launcher/v1"
	runtimev1 "gitlab.com/Geneden/phanes/gen/go/phanes/runtime/v1"
	savev1 "gitlab.com/Geneden/phanes/gen/go/phanes/save/v1"
)

func TestGeneratedContractsCompile(t *testing.T) {
	manifest := &cachev1.CacheManifest{
		CacheId: "compile-check",
		SchemaVersion: &commonv1.Version{
			Major: 0,
			Minor: 1,
			Patch: 0,
		},
	}

	buildResult := &builderv1.BuildResult{Success: true, Manifest: manifest}
	runtimeStatus := &runtimev1.RuntimeStatus{State: runtimev1.RuntimeState_RUNTIME_STATE_READY, CacheId: manifest.GetCacheId()}
	launcherStatus := &launcherv1.RuntimeProcessStatus{RuntimeStatus: runtimeStatus}
	profile := &savev1.Profile{ProfileId: "local"}

	if !buildResult.GetSuccess() || launcherStatus.GetRuntimeStatus().GetCacheId() != "compile-check" || profile.GetProfileId() != "local" {
		t.Fatal("generated contract compile check produced unexpected values")
	}
}
