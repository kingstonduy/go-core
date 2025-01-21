# Makefile for Git tag operations

# Default tag (can be overridden by passing TAG=your-tag)
TAG = v0.0.9

# Remote name (usually 'origin', but can be customized)
REMOTE = origin

# Add a new tag and push it to the remote repository
add-tag:
	@echo "Creating new tag: $(TAG)"
	@git tag $(TAG)
	@git push $(REMOTE) $(TAG)
	@echo "Tag $(TAG) created and pushed to $(REMOTE)"

# Delete a tag locally and remotely
delete-tag:
	@echo "Deleting tag: $(TAG)"
	@git tag -d $(TAG)  # Delete the tag locally
	@git push $(REMOTE) :refs/tags/$(TAG)  # Delete the tag remotely
	@echo "Tag $(TAG) deleted locally and remotely"

# Help message (optional)
help:
	@echo "Makefile for Git tag operations"
	@echo "Usage:"
	@echo "  make add-tag     # Add a new tag and push it"
	@echo "  make delete-tag  # Delete a tag locally and remotely"
	@echo "  make help        # Show this help message"
