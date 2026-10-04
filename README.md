# ROS: Relation Optimizer System

> **A Sequential Bayesian Low-Pass Filter for Personal Bandwidth Defense.**

ROS is a command-line signal detection engine designed to protect finite human attention, energy, and capital from exploiters, bad-faith actors, and cognitive overload.

Treating social interactions not through emotional heuristics or moral judgments, but strictly as a **Sequential Probability Ratio Test (SPRT)**, ROS applies Alan Turing's Banburismus (Deciban additive log-odds) and dual-timescale exponential decay to model relational reliability objectively.

---

## Core Philosophy: Bandwidth Defense, Not Moral Judgment

ROS is **not** a social credit scoring tool or an ethical evaluation platform. It is a strictly private, client-side low-pass filter built on three foundational principles:

1. **Protocol Over Personality**: Results and environmental fortunes are noisy. ROS only measures controllable behavioral decisions and adherence to commitment protocols.
2. **Orthogonal Sampling**: To prevent overconfidence and score inflation, every event is mapped to a single dominant invariant factor rather than correlated multi-counting.
3. **Decay-Driven Forgiveness**: Friction and habitual noise decay exponentially over time. Human relationships should not require manual ledger updates; superficial errors fade automatically, while catastrophic breaches remain permanent.

---

## Mathematical Architecture

### 1. Banburismus & Deciban Scaling
Instead of computing non-linear probability multiplications, ROS operates in **Decibels (dB)** of evidence:

$$\text{Weight of Evidence (dB)} = 10 \log_{10} \frac{P(E \mid H_1)}{P(E \mid H_0)}$$

* **$H_1$**: The hypothesis that the target node reliably operates according to the baseline protocol.
* **$H_0$**: The null hypothesis of typical opportunistic or noisy behavior.

Evidence updates are purely additive:

$$\text{Odds}_{\text{new}} (\text{dB}) = \text{Odds}_{\text{old}} (\text{dB}) + \Delta \text{dB}_{\text{eff}}$$

### 2. Dual-Timescale Exponential Decay
Human memory either over-penalizes past trivialities or forgets critical warnings. ROS formalizes time through structured half-lives ($T_{1/2}$):

$$\Delta \text{dB}_{\text{eff}}(t) = \Delta \text{dB}_{\text{base}} \times 2^{-\frac{\Delta t}{T_{1/2}}}$$

| Inevitability Layer | Half-Life ($T_{1/2}$) | Purpose |
|---|---|---|
| `transient_noise` | **7 Days** | Flattery, one-off minor irritations, mood swings (rapid evaporation). |
| `daily` | **30 Days** | Punctuality, response discipline, routine commitment reliability. |
| `critical_moment` | **$\infty$ (Immutable)** | Crucible choices: scapegoating, contract breaches, or self-sacrificing integrity. |

### 3. Decision Boundaries & Passive Monitoring
Operating under Wald's SPRT framework:
* **$\ge +15.0\text{ dB}$ (`PRIORITY_INVEST`)**: High-confidence cooperative node. Prioritize active communication and investment.
* **$-10.0\text{ dB} < \text{Score} < +15.0\text{ dB}$ (`OBSERVE`)**: Standard evaluation band.
* **$\le -10.0\text{ dB}$ (`EARLY_EXIT`)**: The system automatically transitions the node to `monitoring_only`. Active resource allocation ceases immediately; incoming signals are recorded purely passively without downstream engagement.

---

## Installation & Build

ROS is written in pure Go (CGO-free) using `modernc.org/sqlite`.

```bash
git clone [https://github.com/your-username/ros.git](https://github.com/your-username/ros.git)
cd ros

# Build the standalone binary
make

```

The optimized binary will be compiled to `./bin/ros`.

---

## Usage

### 1. Interactive Console (Recommended)

Launch the zero-cognitive-overhead interactive prompt:

```bash
./bin/ros
# or
./bin/ros i

```

Select numbered options to inspect the dashboard, register nodes, or append fact-based observations.

### 2. CLI Scripting & Automation

#### Register a Node

```bash
./bin/ros add -name "Alice Smith" -type adult_customer
./bin/ros add -name "Bob Unreliable" -type adult_non_customer

```

Supported Persona Types:

* `adult_customer` (Contractual / paying counterparty)
* `adult_non_customer` (Acquaintance / external social connection)
* `employee` (Partner / internal team member)
* `child_customer` (Minor in structured / instructional setting)
* `child_non_customer` (Minor baseline tolerance)

#### Record an Observation

```bash
./bin/ros log -id p_a1b2 -memo "Delayed milestone delivery by 3 days without notice" -sign neg -inev daily
./bin/ros log -id p_c3d4 -memo "Covered unforeseen partner loss out-of-pocket" -sign pos -inev critical_moment

```

#### View Dashboard

```bash
./bin/ros ls

```

Output:

```text
ID         NAME             TYPE                 STATUS           SCORE(dB)  EVENTS   VERDICT        
------------------------------------------------------------------------------------------------------
p_c3d4     Charlie Loyal    employee             active             +20.00   1        PRIORITY_INVEST
p_a1b2     Alice Smith      adult_customer       active              +4.00   2        OBSERVE        
p_bad1     Bob Unreliable   adult_non_customer   monitoring_only    -12.00   3        EARLY_EXIT     

```

---

## Testing

Verify the mathematical decay models, threshold transitions, and database views using in-memory SQLite:

```bash
go test -v ./internal/usecase/...

```

---

## Security & Privacy Notice

* **Local-First**: All data is stored in a single, local SQLite database (`ros.db`).
* **Zero Telemetry**: No network calls, telemetry, or external analytics exist in the codebase.
* **Confidentiality Warning**: If you fork or publish your repository, ensure `ros.db` remains inside `.gitignore`. Never commit factual observation logs regarding actual counterparties.

---

## License

MIT License. Designed for individuals who value objective clarity and cognitive boundaries.
