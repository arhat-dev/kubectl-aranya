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

#
# linux
#
package.kubectl-aranya.deb.amd64:
	sh scripts/package/package.sh $@

package.kubectl-aranya.deb.armv6:
	sh scripts/package/package.sh $@

package.kubectl-aranya.deb.armv7:
	sh scripts/package/package.sh $@

package.kubectl-aranya.deb.arm64:
	sh scripts/package/package.sh $@

package.kubectl-aranya.deb.all: \
	package.kubectl-aranya.deb.amd64 \
	package.kubectl-aranya.deb.armv6 \
	package.kubectl-aranya.deb.armv7 \
	package.kubectl-aranya.deb.arm64

package.kubectl-aranya.rpm.amd64:
	sh scripts/package/package.sh $@

package.kubectl-aranya.rpm.armv7:
	sh scripts/package/package.sh $@

package.kubectl-aranya.rpm.arm64:
	sh scripts/package/package.sh $@

package.kubectl-aranya.rpm.all: \
	package.kubectl-aranya.rpm.amd64 \
	package.kubectl-aranya.rpm.armv7 \
	package.kubectl-aranya.rpm.arm64

package.kubectl-aranya.linux.all: \
	package.kubectl-aranya.deb.all \
	package.kubectl-aranya.rpm.all

#
# windows
#

package.kubectl-aranya.msi.amd64:
	sh scripts/package/package.sh $@

package.kubectl-aranya.msi.arm64:
	sh scripts/package/package.sh $@

package.kubectl-aranya.msi.all: \
	package.kubectl-aranya.msi.amd64 \
	package.kubectl-aranya.msi.arm64

package.kubectl-aranya.windows.all: \
	package.kubectl-aranya.msi.all

#
# darwin
#

package.kubectl-aranya.pkg.amd64:
	sh scripts/package/package.sh $@

package.kubectl-aranya.pkg.arm64:
	sh scripts/package/package.sh $@

package.kubectl-aranya.pkg.all: \
	package.kubectl-aranya.pkg.amd64 \
	package.kubectl-aranya.pkg.arm64

package.kubectl-aranya.darwin.all: \
	package.kubectl-aranya.pkg.all
