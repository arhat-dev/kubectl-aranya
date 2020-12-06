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
package.kubectl-aranya-pf.deb.amd64:
	sh scripts/package/package.sh $@

package.kubectl-aranya-pf.deb.armv6:
	sh scripts/package/package.sh $@

package.kubectl-aranya-pf.deb.armv7:
	sh scripts/package/package.sh $@

package.kubectl-aranya-pf.deb.arm64:
	sh scripts/package/package.sh $@

package.kubectl-aranya-pf.deb.all: \
	package.kubectl-aranya-pf.deb.amd64 \
	package.kubectl-aranya-pf.deb.armv6 \
	package.kubectl-aranya-pf.deb.armv7 \
	package.kubectl-aranya-pf.deb.arm64

package.kubectl-aranya-pf.rpm.amd64:
	sh scripts/package/package.sh $@

package.kubectl-aranya-pf.rpm.armv7:
	sh scripts/package/package.sh $@

package.kubectl-aranya-pf.rpm.arm64:
	sh scripts/package/package.sh $@

package.kubectl-aranya-pf.rpm.all: \
	package.kubectl-aranya-pf.rpm.amd64 \
	package.kubectl-aranya-pf.rpm.armv7 \
	package.kubectl-aranya-pf.rpm.arm64

package.kubectl-aranya-pf.linux.all: \
	package.kubectl-aranya-pf.deb.all \
	package.kubectl-aranya-pf.rpm.all

#
# windows
#

package.kubectl-aranya-pf.msi.amd64:
	sh scripts/package/package.sh $@

package.kubectl-aranya-pf.msi.arm64:
	sh scripts/package/package.sh $@

package.kubectl-aranya-pf.msi.all: \
	package.kubectl-aranya-pf.msi.amd64 \
	package.kubectl-aranya-pf.msi.arm64

package.kubectl-aranya-pf.windows.all: \
	package.kubectl-aranya-pf.msi.all

#
# darwin
#

package.kubectl-aranya-pf.pkg.amd64:
	sh scripts/package/package.sh $@

package.kubectl-aranya-pf.pkg.arm64:
	sh scripts/package/package.sh $@

package.kubectl-aranya-pf.pkg.all: \
	package.kubectl-aranya-pf.pkg.amd64 \
	package.kubectl-aranya-pf.pkg.arm64

package.kubectl-aranya-pf.darwin.all: \
	package.kubectl-aranya-pf.pkg.all
