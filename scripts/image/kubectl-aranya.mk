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

# build
image.build.kubectl-aranya.linux.x86:
	sh scripts/image/build.sh $@

image.build.kubectl-aranya.linux.amd64:
	sh scripts/image/build.sh $@

image.build.kubectl-aranya.linux.armv5:
	sh scripts/image/build.sh $@

image.build.kubectl-aranya.linux.armv6:
	sh scripts/image/build.sh $@

image.build.kubectl-aranya.linux.armv7:
	sh scripts/image/build.sh $@

image.build.kubectl-aranya.linux.arm64:
	sh scripts/image/build.sh $@

image.build.kubectl-aranya.linux.ppc64le:
	sh scripts/image/build.sh $@

image.build.kubectl-aranya.linux.mips64le:
	sh scripts/image/build.sh $@

image.build.kubectl-aranya.linux.s390x:
	sh scripts/image/build.sh $@

image.build.kubectl-aranya.linux.all: \
	image.build.kubectl-aranya.linux.amd64 \
	image.build.kubectl-aranya.linux.arm64 \
	image.build.kubectl-aranya.linux.armv7 \
	image.build.kubectl-aranya.linux.armv6 \
	image.build.kubectl-aranya.linux.armv5 \
	image.build.kubectl-aranya.linux.x86 \
	image.build.kubectl-aranya.linux.s390x \
	image.build.kubectl-aranya.linux.ppc64le \
	image.build.kubectl-aranya.linux.mips64le

image.build.kubectl-aranya.windows.amd64:
	sh scripts/image/build.sh $@

image.build.kubectl-aranya.windows.armv7:
	sh scripts/image/build.sh $@

image.build.kubectl-aranya.windows.all: \
	image.build.kubectl-aranya.windows.amd64 \
	image.build.kubectl-aranya.windows.armv7

# push
image.push.kubectl-aranya.linux.x86:
	sh scripts/image/push.sh $@

image.push.kubectl-aranya.linux.amd64:
	sh scripts/image/push.sh $@

image.push.kubectl-aranya.linux.armv5:
	sh scripts/image/push.sh $@

image.push.kubectl-aranya.linux.armv6:
	sh scripts/image/push.sh $@

image.push.kubectl-aranya.linux.armv7:
	sh scripts/image/push.sh $@

image.push.kubectl-aranya.linux.arm64:
	sh scripts/image/push.sh $@

image.push.kubectl-aranya.linux.ppc64le:
	sh scripts/image/push.sh $@

image.push.kubectl-aranya.linux.mips64le:
	sh scripts/image/push.sh $@

image.push.kubectl-aranya.linux.s390x:
	sh scripts/image/push.sh $@

image.push.kubectl-aranya.linux.all: \
	image.push.kubectl-aranya.linux.amd64 \
	image.push.kubectl-aranya.linux.arm64 \
	image.push.kubectl-aranya.linux.armv7 \
	image.push.kubectl-aranya.linux.armv6 \
	image.push.kubectl-aranya.linux.armv5 \
	image.push.kubectl-aranya.linux.x86 \
	image.push.kubectl-aranya.linux.s390x \
	image.push.kubectl-aranya.linux.ppc64le \
	image.push.kubectl-aranya.linux.mips64le

image.push.kubectl-aranya.windows.amd64:
	sh scripts/image/push.sh $@

image.push.kubectl-aranya.windows.armv7:
	sh scripts/image/push.sh $@

image.push.kubectl-aranya.windows.all: \
	image.push.kubectl-aranya.windows.amd64 \
	image.push.kubectl-aranya.windows.armv7
