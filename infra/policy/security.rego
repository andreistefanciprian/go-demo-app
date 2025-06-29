package main

# Deny if any container doesn't have runAsNonRoot set to true
deny contains msg if {
    input.kind == "Deployment"
    container := input.spec.template.spec.containers[_]
    not container.securityContext.runAsNonRoot == true
    msg := sprintf("Container '%s' must have runAsNonRoot set to true", [container.name])
}

# Deny if pod security context allows root user
deny contains msg if {
    input.kind == "Deployment"
    input.spec.template.spec.securityContext.runAsUser == 0
    msg := "Pod securityContext cannot run as root user (UID 0)"
}

# Deny if container security context allows root user
deny contains msg if {
    input.kind == "Deployment"
    container := input.spec.template.spec.containers[_]
    container.securityContext.runAsUser == 0
    msg := sprintf("Container '%s' cannot run as root user (UID 0)", [container.name])
}

# Deny if runAsNonRoot is not set at pod level
deny contains msg if {
    input.kind == "Deployment"
    not input.spec.template.spec.securityContext.runAsNonRoot == true
    msg := "Pod securityContext must have runAsNonRoot set to true"
}

# Deny if capabilities are not dropped
deny contains msg if {
    input.kind == "Deployment"
    container := input.spec.template.spec.containers[_]
    not container.securityContext.capabilities.drop
    msg := sprintf("Container '%s' must drop all capabilities", [container.name])
}

# Deny if ALL capabilities are not dropped
deny contains msg if {
    input.kind == "Deployment"
    container := input.spec.template.spec.containers[_]
    capabilities := container.securityContext.capabilities.drop
    not "ALL" in capabilities
    msg := sprintf("Container '%s' must drop ALL capabilities", [container.name])
}

# Deny if allowPrivilegeEscalation is not false
deny contains msg if {
    input.kind == "Deployment"
    container := input.spec.template.spec.containers[_]
    not container.securityContext.allowPrivilegeEscalation == false
    msg := sprintf("Container '%s' must set allowPrivilegeEscalation to false", [container.name])
}

# Deny if readOnlyRootFilesystem is not true
deny contains msg if {
    input.kind == "Deployment"
    container := input.spec.template.spec.containers[_]
    not container.securityContext.readOnlyRootFilesystem == true
    msg := sprintf("Container '%s' must have readOnlyRootFilesystem set to true", [container.name])
}
