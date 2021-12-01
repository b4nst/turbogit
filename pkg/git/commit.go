package git

import git "github.com/libgit2/git2go/v33"

// RepoTree return the current index tree
func RepoTree(r *git.Repository) (*git.Tree, error) {
	idx, err := r.Index()
	if err != nil {
		return nil, err
	}
	treeId, err := idx.WriteTree()
	if err != nil {
		return nil, err
	}
	tree, err := r.LookupTree(treeId)
	if err != nil {
		return nil, err
	}
	return tree, nil
}

// Commit creates a new commit with the current tree.
func Commit(r *git.Repository, msg string) (*git.Commit, error) {
	// Signature
	sig, err := r.DefaultSignature()
	if err != nil {
		return nil, err
	}
	// Tree
	tree, err := RepoTree(r)
	if err != nil {
		return nil, err
	}
	// Parents
	parents := []*git.Commit{}
	head, err := r.Head()
	if err == nil { // We found head
		headRef, err := r.LookupCommit(head.Target())
		if err != nil {
			return nil, err
		}
		parents = append(parents, headRef)
	}

	oid, err := r.CreateCommit("HEAD", sig, sig, msg, tree, parents...)
	if err != nil {
		return nil, err
	}
	return r.LookupCommit(oid)
}

// Amend amends the HEAD commit
func Amend(ca *git.Commit, msg string) (*git.Commit, error) {
	r := ca.Object.Owner()
	// Signature
	sig, err := r.DefaultSignature()
	if err != nil {
		return nil, err
	}
	// Tree
	tree, err := RepoTree(r)
	if err != nil {
		return nil, err
	}
	oid, err := ca.Amend("HEAD", ca.Author(), sig, msg, tree)
	if err != nil {
		return nil, err
	}
	return r.LookupCommit(oid)
}
