srpm:
	dnf install -y git
	curl -LOf https://src.fedoraproject.org/rpms/butane/raw/rawhide/f/butane.spec
	version=$$(git describe --always --tags | sed -e 's,-,\.,g' -e 's,^v,,'); \
	git archive --format=tar --prefix=butane-$$version/ HEAD | gzip > butane-$$version.tar.gz; \
	sed -ie "s,^Version:.*,Version: $$version," butane.spec
	sed -ie 's/^Patch/# Patch/g' butane.spec  # we don't want any downstream patches
	rpmbuild -bs --define "_sourcedir ${PWD}" --define "_specdir ${PWD}" --define "_builddir ${PWD}" --define "_srcrpmdir ${PWD}" --define "_rpmdir ${PWD}" --define "_buildrootdir ${PWD}/.build" butane.spec
	mv *.src.rpm $$outdir
