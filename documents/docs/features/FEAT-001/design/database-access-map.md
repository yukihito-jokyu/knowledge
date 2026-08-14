# FEAT-001 Database Access Map（v1）

| Operation | Read | Write | Derived Index | Transaction |
| --- | --- | --- | --- | --- |
| search-text | assertion_lexical_index, assertions, assertion_revisions, revision_scopes, concepts, assertion_concepts | なし | assertion_lexical_index | read-only | 
| search-concept | concepts, concept_terms, assertion_concepts, assertions, assertion_revisions, revision_scopes | なし | idx_concept_terms_concept, idx_assertion_concepts_concept | read-only |
| search-related | relations, assertions, concepts, assertion_revisions | なし | idx_relations_source, idx_relations_target | read-only |
| get | assertions, assertion_revisions, revision_scopes, temporal_metadata, assertion_concepts, concepts, concept_aliases, assertion_aliases | なし | PK/FK | read-only |
| get-evidence | assertions, evidence | なし | idx_evidence_assertion | read-only |
| search-contradictions | relations, assertions, assertion_revisions, concepts, concept_aliases, assertion_concepts | なし | idx_relations_source/target | read-only |
| search-temporal | temporal_metadata, assertions, assertion_revisions, revision_scopes, assertion_concepts, concepts, concept_aliases | なし | idx_temporal_window | read-only |
| create | assertions, assertion_revisions, revision_scopes, evidence, concepts, concept_terms, concept_aliases, assertion_concepts, assertion_aliases, temporal_metadata | assertions, assertion_revisions, revision_scopes, evidence, concepts, concept_terms, concept_aliases, assertion_concepts, assertion_aliases, temporal_metadata, relations, assertion_lexical_index | assertion_lexical_index、scope/concept/alias/relation indexes | 1 write transaction |
| attach-evidence | assertions | evidence | idx_evidence_assertion | 1 write transaction |
| revise | assertions, assertion_revisions, revision_scopes, temporal_metadata, assertion_concepts, concepts, concept_aliases, assertion_aliases, assertion_lexical_index | assertion_revisions, revision_scopes, temporal_metadata, assertions.current_revision, assertion_lexical_index | assertion_lexical_index、scope indexes | 1 write transaction |
| supersede | assertions, relations | relations | relation indexes | 1 write transaction |

すべての write transaction は validation、base record 確認、write、Index整合、commit の順で行う。どこかで失敗すれば rollback し、成功 response は commit 後だけに返す。
