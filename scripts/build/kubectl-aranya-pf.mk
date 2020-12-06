# Copyright 2020 The arhat.dev Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# native
kubectl-aranya-pf:
	sh scripts/build/build.sh $@

# linux
kubectl-aranya-pf.linux.x86:
	sh scripts/build/build.sh $@

kubectl-aranya-pf.linux.amd64:
	sh scripts/build/build.sh $@

kubectl-aranya-pf.linux.armv5:
	sh scripts/build/build.sh $@

kubectl-aranya-pf.linux.armv6:
	sh scripts/build/build.sh $@

kubectl-aranya-pf.linux.armv7:
	sh scripts/build/build.sh $@

kubectl-aranya-pf.linux.arm64:
	sh scripts/build/build.sh $@

kubectl-aranya-pf.linux.mips:
	sh scripts/build/build.sh $@

kubectl-aranya-pf.linux.mipshf:
	sh scripts/build/build.sh $@

kubectl-aranya-pf.linux.mipsle:
	sh scripts/build/build.sh $@

kubectl-aranya-pf.linux.mipslehf:
	sh scripts/build/build.sh $@

kubectl-aranya-pf.linux.mips64:
	sh scripts/build/build.sh $@

kubectl-aranya-pf.linux.mips64hf:
	sh scripts/build/build.sh $@

kubectl-aranya-pf.linux.mips64le:
	sh scripts/build/build.sh $@

kubectl-aranya-pf.linux.mips64lehf:
	sh scripts/build/build.sh $@

kubectl-aranya-pf.linux.ppc64:
	sh scripts/build/build.sh $@

kubectl-aranya-pf.linux.ppc64le:
	sh scripts/build/build.sh $@

kubectl-aranya-pf.linux.s390x:
	sh scripts/build/build.sh $@

kubectl-aranya-pf.linux.riscv64:
	sh scripts/build/build.sh $@

kubectl-aranya-pf.linux.all: \
	kubectl-aranya-pf.linux.x86 \
	kubectl-aranya-pf.linux.amd64 \
	kubectl-aranya-pf.linux.armv5 \
	kubectl-aranya-pf.linux.armv6 \
	kubectl-aranya-pf.linux.armv7 \
	kubectl-aranya-pf.linux.arm64 \
	kubectl-aranya-pf.linux.mips \
	kubectl-aranya-pf.linux.mipshf \
	kubectl-aranya-pf.linux.mipsle \
	kubectl-aranya-pf.linux.mipslehf \
	kubectl-aranya-pf.linux.mips64 \
	kubectl-aranya-pf.linux.mips64hf \
	kubectl-aranya-pf.linux.mips64le \
	kubectl-aranya-pf.linux.mips64lehf \
	kubectl-aranya-pf.linux.ppc64 \
	kubectl-aranya-pf.linux.ppc64le \
	kubectl-aranya-pf.linux.s390x \
	kubectl-aranya-pf.linux.riscv64

kubectl-aranya-pf.darwin.amd64:
	sh scripts/build/build.sh $@

# # currently darwin/arm64 build will fail due to golang link error
# kubectl-aranya-pf.darwin.arm64:
# 	sh scripts/build/build.sh $@

kubectl-aranya-pf.darwin.all: \
	kubectl-aranya-pf.darwin.amd64

kubectl-aranya-pf.windows.x86:
	sh scripts/build/build.sh $@

kubectl-aranya-pf.windows.amd64:
	sh scripts/build/build.sh $@

kubectl-aranya-pf.windows.armv5:
	sh scripts/build/build.sh $@

kubectl-aranya-pf.windows.armv6:
	sh scripts/build/build.sh $@

kubectl-aranya-pf.windows.armv7:
	sh scripts/build/build.sh $@

# # currently no support for windows/arm64
# kubectl-aranya-pf.windows.arm64:
# 	sh scripts/build/build.sh $@

kubectl-aranya-pf.windows.all: \
	kubectl-aranya-pf.windows.x86 \
	kubectl-aranya-pf.windows.amd64 \
	kubectl-aranya-pf.windows.armv5 \
	kubectl-aranya-pf.windows.armv6 \
	kubectl-aranya-pf.windows.armv7

# # android build requires android sdk
# kubectl-aranya-pf.android.amd64:
# 	sh scripts/build/build.sh $@

# kubectl-aranya-pf.android.x86:
# 	sh scripts/build/build.sh $@

# kubectl-aranya-pf.android.armv5:
# 	sh scripts/build/build.sh $@

# kubectl-aranya-pf.android.armv6:
# 	sh scripts/build/build.sh $@

# kubectl-aranya-pf.android.armv7:
# 	sh scripts/build/build.sh $@

# kubectl-aranya-pf.android.arm64:
# 	sh scripts/build/build.sh $@

# kubectl-aranya-pf.android.all: \
# 	kubectl-aranya-pf.android.amd64 \
# 	kubectl-aranya-pf.android.arm64 \
# 	kubectl-aranya-pf.android.x86 \
# 	kubectl-aranya-pf.android.armv7 \
# 	kubectl-aranya-pf.android.armv5 \
# 	kubectl-aranya-pf.android.armv6

kubectl-aranya-pf.freebsd.amd64:
	sh scripts/build/build.sh $@

kubectl-aranya-pf.freebsd.x86:
	sh scripts/build/build.sh $@

kubectl-aranya-pf.freebsd.armv5:
	sh scripts/build/build.sh $@

kubectl-aranya-pf.freebsd.armv6:
	sh scripts/build/build.sh $@

kubectl-aranya-pf.freebsd.armv7:
	sh scripts/build/build.sh $@

kubectl-aranya-pf.freebsd.arm64:
	sh scripts/build/build.sh $@

kubectl-aranya-pf.freebsd.all: \
	kubectl-aranya-pf.freebsd.amd64 \
	kubectl-aranya-pf.freebsd.arm64 \
	kubectl-aranya-pf.freebsd.armv7 \
	kubectl-aranya-pf.freebsd.x86 \
	kubectl-aranya-pf.freebsd.armv5 \
	kubectl-aranya-pf.freebsd.armv6

kubectl-aranya-pf.netbsd.amd64:
	sh scripts/build/build.sh $@

kubectl-aranya-pf.netbsd.x86:
	sh scripts/build/build.sh $@

kubectl-aranya-pf.netbsd.armv5:
	sh scripts/build/build.sh $@

kubectl-aranya-pf.netbsd.armv6:
	sh scripts/build/build.sh $@

kubectl-aranya-pf.netbsd.armv7:
	sh scripts/build/build.sh $@

kubectl-aranya-pf.netbsd.arm64:
	sh scripts/build/build.sh $@

kubectl-aranya-pf.netbsd.all: \
	kubectl-aranya-pf.netbsd.amd64 \
	kubectl-aranya-pf.netbsd.arm64 \
	kubectl-aranya-pf.netbsd.armv7 \
	kubectl-aranya-pf.netbsd.x86 \
	kubectl-aranya-pf.netbsd.armv5 \
	kubectl-aranya-pf.netbsd.armv6

kubectl-aranya-pf.openbsd.amd64:
	sh scripts/build/build.sh $@

kubectl-aranya-pf.openbsd.x86:
	sh scripts/build/build.sh $@

kubectl-aranya-pf.openbsd.armv5:
	sh scripts/build/build.sh $@

kubectl-aranya-pf.openbsd.armv6:
	sh scripts/build/build.sh $@

kubectl-aranya-pf.openbsd.armv7:
	sh scripts/build/build.sh $@

kubectl-aranya-pf.openbsd.arm64:
	sh scripts/build/build.sh $@

kubectl-aranya-pf.openbsd.all: \
	kubectl-aranya-pf.openbsd.amd64 \
	kubectl-aranya-pf.openbsd.arm64 \
	kubectl-aranya-pf.openbsd.armv7 \
	kubectl-aranya-pf.openbsd.x86 \
	kubectl-aranya-pf.openbsd.armv5 \
	kubectl-aranya-pf.openbsd.armv6

kubectl-aranya-pf.solaris.amd64:
	sh scripts/build/build.sh $@

kubectl-aranya-pf.aix.ppc64:
	sh scripts/build/build.sh $@

kubectl-aranya-pf.dragonfly.amd64:
	sh scripts/build/build.sh $@

kubectl-aranya-pf.plan9.amd64:
	sh scripts/build/build.sh $@

kubectl-aranya-pf.plan9.x86:
	sh scripts/build/build.sh $@

kubectl-aranya-pf.plan9.armv5:
	sh scripts/build/build.sh $@

kubectl-aranya-pf.plan9.armv6:
	sh scripts/build/build.sh $@

kubectl-aranya-pf.plan9.armv7:
	sh scripts/build/build.sh $@

kubectl-aranya-pf.plan9.all: \
	kubectl-aranya-pf.plan9.amd64 \
	kubectl-aranya-pf.plan9.armv7 \
	kubectl-aranya-pf.plan9.x86 \
	kubectl-aranya-pf.plan9.armv5 \
	kubectl-aranya-pf.plan9.armv6
