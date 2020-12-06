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
kubectl-aranya:
	sh scripts/build/build.sh $@

# linux
kubectl-aranya.linux.x86:
	sh scripts/build/build.sh $@

kubectl-aranya.linux.amd64:
	sh scripts/build/build.sh $@

kubectl-aranya.linux.armv5:
	sh scripts/build/build.sh $@

kubectl-aranya.linux.armv6:
	sh scripts/build/build.sh $@

kubectl-aranya.linux.armv7:
	sh scripts/build/build.sh $@

kubectl-aranya.linux.arm64:
	sh scripts/build/build.sh $@

kubectl-aranya.linux.mips:
	sh scripts/build/build.sh $@

kubectl-aranya.linux.mipshf:
	sh scripts/build/build.sh $@

kubectl-aranya.linux.mipsle:
	sh scripts/build/build.sh $@

kubectl-aranya.linux.mipslehf:
	sh scripts/build/build.sh $@

kubectl-aranya.linux.mips64:
	sh scripts/build/build.sh $@

kubectl-aranya.linux.mips64hf:
	sh scripts/build/build.sh $@

kubectl-aranya.linux.mips64le:
	sh scripts/build/build.sh $@

kubectl-aranya.linux.mips64lehf:
	sh scripts/build/build.sh $@

kubectl-aranya.linux.ppc64:
	sh scripts/build/build.sh $@

kubectl-aranya.linux.ppc64le:
	sh scripts/build/build.sh $@

kubectl-aranya.linux.s390x:
	sh scripts/build/build.sh $@

kubectl-aranya.linux.riscv64:
	sh scripts/build/build.sh $@

kubectl-aranya.linux.all: \
	kubectl-aranya.linux.x86 \
	kubectl-aranya.linux.amd64 \
	kubectl-aranya.linux.armv5 \
	kubectl-aranya.linux.armv6 \
	kubectl-aranya.linux.armv7 \
	kubectl-aranya.linux.arm64 \
	kubectl-aranya.linux.mips \
	kubectl-aranya.linux.mipshf \
	kubectl-aranya.linux.mipsle \
	kubectl-aranya.linux.mipslehf \
	kubectl-aranya.linux.mips64 \
	kubectl-aranya.linux.mips64hf \
	kubectl-aranya.linux.mips64le \
	kubectl-aranya.linux.mips64lehf \
	kubectl-aranya.linux.ppc64 \
	kubectl-aranya.linux.ppc64le \
	kubectl-aranya.linux.s390x \
	kubectl-aranya.linux.riscv64

kubectl-aranya.darwin.amd64:
	sh scripts/build/build.sh $@

# # currently darwin/arm64 build will fail due to golang link error
# kubectl-aranya.darwin.arm64:
# 	sh scripts/build/build.sh $@

kubectl-aranya.darwin.all: \
	kubectl-aranya.darwin.amd64

kubectl-aranya.windows.x86:
	sh scripts/build/build.sh $@

kubectl-aranya.windows.amd64:
	sh scripts/build/build.sh $@

kubectl-aranya.windows.armv5:
	sh scripts/build/build.sh $@

kubectl-aranya.windows.armv6:
	sh scripts/build/build.sh $@

kubectl-aranya.windows.armv7:
	sh scripts/build/build.sh $@

# # currently no support for windows/arm64
# kubectl-aranya.windows.arm64:
# 	sh scripts/build/build.sh $@

kubectl-aranya.windows.all: \
	kubectl-aranya.windows.x86 \
	kubectl-aranya.windows.amd64 \
	kubectl-aranya.windows.armv5 \
	kubectl-aranya.windows.armv6 \
	kubectl-aranya.windows.armv7

# # android build requires android sdk
# kubectl-aranya.android.amd64:
# 	sh scripts/build/build.sh $@

# kubectl-aranya.android.x86:
# 	sh scripts/build/build.sh $@

# kubectl-aranya.android.armv5:
# 	sh scripts/build/build.sh $@

# kubectl-aranya.android.armv6:
# 	sh scripts/build/build.sh $@

# kubectl-aranya.android.armv7:
# 	sh scripts/build/build.sh $@

# kubectl-aranya.android.arm64:
# 	sh scripts/build/build.sh $@

# kubectl-aranya.android.all: \
# 	kubectl-aranya.android.amd64 \
# 	kubectl-aranya.android.arm64 \
# 	kubectl-aranya.android.x86 \
# 	kubectl-aranya.android.armv7 \
# 	kubectl-aranya.android.armv5 \
# 	kubectl-aranya.android.armv6

kubectl-aranya.freebsd.amd64:
	sh scripts/build/build.sh $@

kubectl-aranya.freebsd.x86:
	sh scripts/build/build.sh $@

kubectl-aranya.freebsd.armv5:
	sh scripts/build/build.sh $@

kubectl-aranya.freebsd.armv6:
	sh scripts/build/build.sh $@

kubectl-aranya.freebsd.armv7:
	sh scripts/build/build.sh $@

kubectl-aranya.freebsd.arm64:
	sh scripts/build/build.sh $@

kubectl-aranya.freebsd.all: \
	kubectl-aranya.freebsd.amd64 \
	kubectl-aranya.freebsd.arm64 \
	kubectl-aranya.freebsd.armv7 \
	kubectl-aranya.freebsd.x86 \
	kubectl-aranya.freebsd.armv5 \
	kubectl-aranya.freebsd.armv6

kubectl-aranya.netbsd.amd64:
	sh scripts/build/build.sh $@

kubectl-aranya.netbsd.x86:
	sh scripts/build/build.sh $@

kubectl-aranya.netbsd.armv5:
	sh scripts/build/build.sh $@

kubectl-aranya.netbsd.armv6:
	sh scripts/build/build.sh $@

kubectl-aranya.netbsd.armv7:
	sh scripts/build/build.sh $@

kubectl-aranya.netbsd.arm64:
	sh scripts/build/build.sh $@

kubectl-aranya.netbsd.all: \
	kubectl-aranya.netbsd.amd64 \
	kubectl-aranya.netbsd.arm64 \
	kubectl-aranya.netbsd.armv7 \
	kubectl-aranya.netbsd.x86 \
	kubectl-aranya.netbsd.armv5 \
	kubectl-aranya.netbsd.armv6

kubectl-aranya.openbsd.amd64:
	sh scripts/build/build.sh $@

kubectl-aranya.openbsd.x86:
	sh scripts/build/build.sh $@

kubectl-aranya.openbsd.armv5:
	sh scripts/build/build.sh $@

kubectl-aranya.openbsd.armv6:
	sh scripts/build/build.sh $@

kubectl-aranya.openbsd.armv7:
	sh scripts/build/build.sh $@

kubectl-aranya.openbsd.arm64:
	sh scripts/build/build.sh $@

kubectl-aranya.openbsd.all: \
	kubectl-aranya.openbsd.amd64 \
	kubectl-aranya.openbsd.arm64 \
	kubectl-aranya.openbsd.armv7 \
	kubectl-aranya.openbsd.x86 \
	kubectl-aranya.openbsd.armv5 \
	kubectl-aranya.openbsd.armv6

kubectl-aranya.solaris.amd64:
	sh scripts/build/build.sh $@

kubectl-aranya.aix.ppc64:
	sh scripts/build/build.sh $@

kubectl-aranya.dragonfly.amd64:
	sh scripts/build/build.sh $@

kubectl-aranya.plan9.amd64:
	sh scripts/build/build.sh $@

kubectl-aranya.plan9.x86:
	sh scripts/build/build.sh $@

kubectl-aranya.plan9.armv5:
	sh scripts/build/build.sh $@

kubectl-aranya.plan9.armv6:
	sh scripts/build/build.sh $@

kubectl-aranya.plan9.armv7:
	sh scripts/build/build.sh $@

kubectl-aranya.plan9.all: \
	kubectl-aranya.plan9.amd64 \
	kubectl-aranya.plan9.armv7 \
	kubectl-aranya.plan9.x86 \
	kubectl-aranya.plan9.armv5 \
	kubectl-aranya.plan9.armv6
