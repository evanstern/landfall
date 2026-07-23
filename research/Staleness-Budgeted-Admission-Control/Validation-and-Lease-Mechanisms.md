---
title: Validation and Lease Mechanisms
aliases: [E5 prior art]
tags: [e5, occ, transactional-memory, leases]
type: note
created: 2026-07-23
updated: 2026-07-23
related: [[Staleness-Budgeted-Admission-Control]]
---

# Validation and Lease Mechanisms

Landfall's E5 — a generation lease taken at admission, with landing validation deciding
Landed / Superseded / Stale — is, per the seed survey, a verbatim mechanism transfer from
concurrency control.

## Optimistic concurrency control (Kung & Robinson 1981)

Transactions run in three phases: **read** (work in a private workspace), **validation** (check
serializability against transactions that committed meanwhile), **write** (apply only if
validation passed, else abort and retry). The bet is that conflicts are rare enough that
validate-after beats lock-before ([[_grounding]] §E5;
[ACM TODS](https://dl.acm.org/doi/10.1145/319566.319567)). The structure maps onto ask → flight
→ land: the oracle call is the read phase, landing is validation, Superseded/Stale are aborts.

## TL2's global version clock (Dice, Shalev & Shavit 2006)

TL2 validates a transaction's read-set against a **global version-clock** at commit time and
stamps written locations with the new clock value — a concrete generation-counter discipline
for detecting "the world moved while you worked" ([[_grounding]] §E5).

## Leases (Gray & Cheriton 1989)

A lease is a **time-bounded grant** of authority over cached state; expiry bounds how stale a
holder's view can be without revalidation. The term and mechanism originate here
([[_grounding]] §E5; [ACM](https://dl.acm.org/doi/10.1145/74851.74870)).

## Adjacent: robotics

The seed survey also notes plan-validity re-checks after planner latency in robotics — validate
that the world still matches the plan's assumptions before executing — as the same landing
check outside databases ([[_grounding]] §E5 via seed survey).

## Grounding

- [[_grounding]] — §E5 (OCC, TL2, leases)
