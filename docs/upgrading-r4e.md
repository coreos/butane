---
title: RHEL for Edge
parent: Upgrading configs
nav_order: 4
---

# Upgrading RHEL for Edge configs

Occasionally, changes are made to RHEL for Edge Butane configs (those that specify `variant: r4e`) that break backward compatibility. While this is not a concern for running machines, since Ignition only runs one time during first boot, it is a concern for those who maintain configuration files. This document serves to detail each of the breaking changes and tries to provide some reasoning for the change. This does not cover all of the changes to the spec - just those that need to be considered when migrating from one version to the next.

{: .no_toc }

1. TOC
{:toc}
