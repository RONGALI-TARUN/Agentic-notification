# Attention Orchestration Platform Architecture

Modern notification systems must *respect user attention* rather than simply broadcasting alerts. Early work by Horvitz et al. described a “notification platform” that sits as an **attention-savvy layer** between incoming messages and the user, weighing the value of information against the cost of interruption. In our design, we take this further: instead of a pure “notification delivery platform,” we build an **Attention Orchestration Platform** that uses AI only in the middle layer. The overall flow becomes:

```
Sources (GitHub, Jira, Webhooks…) 
   → 1. Event Layer 
   → 2. Attention Layer (AI) 
   → 3. Delivery Layer 
   → 4. Feedback Layer 
   → Future Learning
```

Each layer has distinct responsibilities, minimizing where AI is needed and keeping the rest deterministic:

- **Event Layer:** Ingests and normalizes events from all sources (e.g. GitHub PRs, Jira issues, custom webhooks). Here we apply *deterministic logic*: validate fields, deduplicate repeated events, and run data-driven rules to assign a *criticality* or *attention class* (see below).  An **event bus / publish-subscribe** pattern is ideal: each source publishes generic events and subscribers (connectors or processors) consume them without tight coupling. If an event is classified as **Critical**, it bypasses most delays and proceeds immediately. Otherwise it moves to the Attention Layer.

- **Attention Layer (LLM Classification):** This is the AI-controlled heart of the system. Using ML/LLM models we classify and group non-critical events, schedule them, assign priorities, and generate summary digests. For example, Apple’s iOS now uses on-device AI to *summarize* and *prioritize* notifications so users see key alerts (“a group chat is active” or urgent texts). Similarly, our LLM examines each event’s content and context to decide whether it should be sent now, merged with others, delayed until after work hours, or held for a digest.  Critical (urgent) events are kept out of this process – only “standard” and “background” events enter the Attention Layer, simplifying AI decisions.

- **Delivery Layer:** A multi-channel notification hub sends messages to users. In practice this means integrating with systems like **email**, **Slack**, **SMS**, or a notifications API (e.g. [Novu](https://novu.co) or similar). The delivery layer has *no AI* – it merely formats and routes messages. A unified API (as Novu provides) lets us plug in dozens of providers at once. For example, Novu’s multi-channel service lets us send emails, push/SMS, chat messages, etc., through one platform. This keeps delivery reliable and consistent, regardless of upstream logic. The Delivery Layer also observes success/failure and respects any *user preferences* (channels/timezones), which users can set in advance.

- **Feedback Layer:** After delivery, we log user responses (opened/read, dismissed, snoozed, replies, etc.) and allow one-click feedback. Rather than quietly “suppressing” notifications, we present an interactive digest where users can undo or change decisions. For example, a daily summary email or Slack message could show “Merged 8 low-priority Jira comments because you usually ignore them” or “Delayed 3 PR reviews (outside work hours)” and offer buttons: *Always Notify*, *Never Notify*, *Undo*. This interface gathers explicit preferences and trust. It maps to common notification UI patterns: MagicBell’s design guide emphasizes letting users mute, set importance, or mark items as read. By capturing reactions (e.g. user pressed “Never Notify this type”), we feed that back to adjust future rules and model thresholds.

Critically, **all decision details are logged** (not just blindly recomputed). We extend the event schema to include fields like:

```json
{
  "attention_class": "critical | standard | background",
  "decision_state": {
    "status": "sent | merged | digest | delayed | suppressed",
    "reason": "...",           
    "model_version": "...",
    "confidence": 0.94
  }
}
```

This ensures that every notification’s fate is recorded (with reason and model info), so we can always trace “why was this notification hidden six months ago?” without re-running the model.

## Notification Classes: Critical, Standard, Background

To simplify logic, we pre-classify events into three priority tiers. Each tier carries fixed policies:

- **CRITICAL**: Must be delivered immediately (bypassing quiet modes, focus filters, or summarization). Never delayed, merged, or suppressed.  *(Example: a P1 incident alert, security breach, or urgent error.)*

- **STANDARD**: Normal priority. The Attention Layer *may* group, delay, or summarize these if needed. *(Example: a PR review request, new assignment, or team comment.)*

- **BACKGROUND**: Low priority. These events may *not* generate real-time alerts; they can be shown only in digests or dropped altogether if things get busy. *(Example: non-critical GitHub events like “someone starred your repo” or minor comments.)*

This three-tier scheme mirrors mobile interruption levels. For instance, iOS defines **Active** (standard) and **Time-Sensitive** (urgent) notifications that respect focus modes, **Passive** (low-priority, silent) that never interrupt, and **Critical** alerts (which bypass everything). By mapping our tiers accordingly, we hard-wire half the decision: a Critical event is handled before the AI sees it (always sent), and Background events are constrained (often held for digest or filtered out). Standard events are the ones the LLM primarily judges.

## Attention Budget vs Suppression

We avoid the term “suppression” because users don’t like feeling ignored. Instead, we think in terms of an **attention budget**. Every person has limited slots for alerts each day. As one recent analysis notes, a user’s *notification budget* is only a few actionable updates per day. In fact, smartphone users already endure 50–60 pushes daily; adding even a handful more from an AI assistant can lead to churn. Therefore, we treat each sent notification as a withdrawal from a finite account. Once the budget is spent, no further “standard” alerts go out — only true critical items can break through.

This reframing is powerful: rather than secretly muting a notification (“suppressed”), we simply allocate from the limited budget. Critical messages cost **zero** attention budget (they must go through). Standard messages *consume budget*. Background messages cost nothing because they don’t draw your immediate attention (they live in the digest or vanish). When the budget is exceeded, we pause sending Standard alerts until it resets the next day. This honors user attention by design, aligning with UX research that finds users severely limit themselves to only a few ‘interruptions’ daily.

## The Central Role of Digests

Because attentional bandwidth is tight, **digests** are a crucial feature for transparency and trust. A digest gathers background and delayed items into a single summary.  Novu describes a notification digest as “one notification instead of multiple separate messages” when a user would otherwise be bombarded. For example, instead of firing off 100 emails for 100 comments, a daily digest could say “8 new comments from your team.” Empirical guides show that digests dramatically reduce overload: “It solves the problem of users getting too many notifications at once,” Novu explains, and “digest notifications work by collecting multiple trigger events, aggregating them into a single message”.

In our product, the digest becomes the *explainability interface*. Rather than users seeing “AI decided something”, they see: 

> **Today’s Digest:**  
> ✓ *Merged* 8 Jira comments (you usually ignore these).  
> ✓ *Delayed* 3 GitHub PR reviews (outside work hours).  
> ✓ *Sent Immediately* – P1 incident, Security alert.  
> ✓ *Hidden* – 12 repo stars.  
> *[Undo] [Always Notify] [Never Notify]* buttons on each item.

This kind of feedback window does double duty. It exposes the *reasoning* behind every move (“because you usually ignore it” or “because it’s not urgent”), building user trust and giving them control. By clicking *Always Notify*, a user effectively adjusts the model or rule for that event type moving forward; *Never Notify* does the opposite. This satisfies product-lead principles of letting users set importance or mute channels (as recommended by design experts) without first building a heavy dashboard. In short, the digest is the explainability and training data center: it collects user corrections that feed back into the system.

## User Feedback and Preferences

Direct user feedback is key to personalization. We embed simple controls in notifications so users can instantly signal preferences. (MagicBell’s UI guidelines similarly advise actions like *“set importance”* or *“mute”* on notification items.) Examples of interactions:

- **Undo / Revert:** allow reversing an AI decision on an alert.  
- **Always/Never Notify:** binary choices that immediately update user policy for that event.  
- **Mute Similar:** silence this category (mapped to our “never notify” for that class).  
- **Looks Good / Thumbs Up:** positive reinforcement.

Recording these in our Feedback Layer closes the loop: we store the user’s choice and any corrections, then use them to adjust scores, thresholds, or retrain models. The system also respects static user settings: for example, quiet hours or channel preferences (e.g. only Slack, not email) as managed in a user profile or preferences UI. By combining automated attention modeling with explicit user control, we create a transparent, trustable experience.

## Data Schema (Attention Metadata)

To support auditing, we augment each event record with the decisions around it. In addition to the original `criticality` field, we include fields like `attention_class` (Critical/Standard/Background) and a nested `decision_state` object. This records the final outcome (*sent*, *merged*, *digest*, *delayed*, or *suppressed*), the reason (e.g. “outside office hours”), the model version, and confidence. Example JSON snippet:

```json
{
  "attention_class": "standard",
  "decision_state": {
      "status": "delayed",
      "reason": "low priority, queued for digest",
      "model_version": "v1.2.3",
      "confidence": 0.86
  }
}
```

Storing decisions like this ensures we never have to guess why something was dropped — the exact rationale is logged. It also isolates the explainable logic (a simple rule or threshold comparison) from the LLM output. In other words, the ML only provides classification scores, while the recorded decision (with numeric cost/benefit logic) provides a precise traceable reason.

## Explainability: Machine vs Human Explanation

In line with best practices in explainable AI, we keep the core decision logic **deterministic** and only use the LLM for presentation. For example, we might determine in code that an event was merged because it was from the same entity and actor within 10 minutes (confidence 96%). We then have the LLM rephrase that rule into plain English in the digest. This avoids “hallucinated” explanations and ensures consistency over time (if we swap models or update training data, past rationale remains the same since it’s logged). Think of it as: 

> Machine Logic (exact condition met) → **LLM** → Human-Readable Explanation.

This approach aligns with enterprise trust requirements: the underlying logic can be inspected, audited, and tested independently of the natural-language layer.

## Multi-Agent Recipients

To enable future agent-to-agent workflows, we abstract recipients as a generic interface. A **Recipient** could be a human (via Email/Slack/SMS), **another AI agent**, a webhook, or a custom workflow endpoint. Concretely, we define something like:

```go
type Recipient interface {
    Receive(Notification) error
}
```

The Delivery Layer simply hands the notification to whatever `Recipient` endpoint is configured. Initially that means human channels, but later it could include bots or workflow triggers. This keeps V1 simple while leaving the door open to “routing a notification to an AI assistant” or other extensions.

## Rule Engine and Pipeline

The one part of this system that calls for careful design is the **rule engine** governing event handling. Every new source (GitHub, Jira, PagerDuty, etc.) should plug into the same framework. We make these rules *data-driven* from day one. Instead of hard-coding “if event.type == X then critical” in code, we define rules in JSON/YAML that the engine loads. This resembles business rule management systems; popular libraries (or even AWS EventBridge-like rules) can match events based on conditions and tag them accordingly. For example, a YAML rule might say:

```yaml
- if:
    source: github
    event: "security_alert"
  then:
    set attention_class: critical
    set reason: "security vulnerability"
```

Because we use a generic subscription model (e.g. an event bus as in the Observer pattern), adding a new connector is *connector + rules*, not rewriting core logic. Over time this same rule framework can support per-org policies, user overrides, and even AI-generated recommendations for better rules. In short, invest in a flexible pipeline early: ingestion → rule classification → queuing by priority. Established patterns (publish/subscribe queues, Chain-of-Responsibility filters, etc.) make this robust and scalable.

## Roadmap Priorities

Based on these design principles, the next steps are:

1. **Core Infrastructure:** Build the ingestion pipeline and schema. Implement the rule engine, event bus, and storage. Ensure every event can be tagged with `attention_class` and persisted.
2. **Reliable Delivery:** Integrate with a notification service (e.g. Novu or custom) for email, Slack, SMS. Guarantee critical alerts (CRITICAL class) never fail.
3. **Attention Layer:** Develop AI classifiers and LLM prompts for Standard events. Implement deduplication, grouping, scheduling (back-off digests, as Novu describes).
4. **Feedback Loop:** Build the interactive digest interfaces (email/Slack) with buttons for Undo/Always/Never. Log all user responses to create training labels.
5. **Learning:** Use collected feedback to refine thresholds and personalize (e.g. adjust a user’s daily budget or scoring based on past dismissals).
6. **Advanced Features:** Once the basics work, add an admin dashboard, fine-grained controls, multi-agent routing, analytics, compliance checks, etc.

Each piece depends on the foundation of the rule/pipeline system. As Horvitz’s work shows, the goal is “systems that respect attention, reduce friction, and help people stay focused on what matters”. By architecting layers cleanly and leveraging AI only where it adds value (classification, summarization, prioritization), we can achieve a notification platform that truly orchestrates user attention.

**Sources:** We drew on industry and research best practices. For example, Novu’s design guide explains digests and multi-channel notifications, Apple’s documentation describes AI-powered notification summaries, and Microsoft’s Attention research (Horvitz et al.) details an intelligent Notification Platform. UX patterns from platforms like MagicBell also influenced our UI and rule-engine approach. These inform the recommendations above.