# Concepts

A **Task** is a single executable function or bound method.
Tasks have **inputs** and **outputs**
defined by parameters and results of the functions.

Multiple *interdependent* tasks come together to form a **Flow**.
Flows can have **flow inputs** and **flow outputs**.
Flows are self-contained:
all inputs come from another task or flow inputs,
and all outputs flow into another task or flow outputs.

Multiple *independent* tasks come together to form a **Parallel**.
Parallels *do not* have inputs or outputs.

Given a Flow or a Parallel, the **Scheduler** runs the different tasks,
ensuring that tasks do not run before their dependencies.
