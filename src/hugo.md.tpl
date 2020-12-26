---
title: "{{ .Title }}"
date: {{ .CreatedAt }}
draft: {{ not .Released }}
---

{{ .Content }}