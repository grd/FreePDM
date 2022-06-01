# FreePDM
***Concept Of Design***

## Requirements

The file below contains 4 tabs.
In the following paragraphs are a explanation of every tab.

[Spreadsheet with requirements](FreePDM_04-Requirements.fods)

### Explanation on the fourth Tab

The fourth Tab contain information about how to handle / create Requirements. The pitfalls of writing them including a link with extended information.  
Read this tab before adding Requirements, otherwise we quite busy fixing Requirement errors.
<!-- Maybe i(Jee-Bee) move the fourth Tab to here -->

### Explanation on the third Tab

The third Tab is added to mention all the assumptions done.
These should be tracked as well

### Explanation on the second Tab

The second tab contains all the definitions.
**All definitions start also within the requirements with a capital letter**.
This is done to make aware there is something _special_ with it.
If a requirement contains a word that is part of the definitions but has no capital letter. please change it.

### Explanation on the first Tab

The first tab is the most important.
This tab contain all the requirements.
Every line contain a 

- requirement id
- level
- traced from
- the requirements
- a note(optionally)

For example there is a requirement:
_The System provide capability to trace Revision history_.
This requirement has _requirement ID_:1 and is a system level requirement so _level_ 0.
and another requirement
_Every File Shall have a Revision_.
This requirement has _ID_ 1 as a parent and so it _trace back to_ requirement _ID_ 1.
Defining the (system )level is more difficult.

## Requirements writing

This is a short, point wise summary about the pitfalls of requirements writing.
the Original Document can be found [here](https://spacese.spacegrant.org/uploads/Requirements-Writing/Writing Good Requirements.pdf)

Good Requirement are necessary, Verifiable, Attainable

How to test:

- **Need** If there is a doubt about the necessity of a requirement, then ask: What is the worst thing that could happen if this requirement were not included? If you do not find an answer of any consequence, then you probably do not need the requirement.
- **Verification** As you write a requirement, determine how you will verify it. Determine the criteria for acceptance. This step will help insure that the requirement is verifiable.
- **Attainable** To be attainable, the requirement must be technically feasible and fit within budget, schedule, and other constraints. If you are uncertain about whether a requirement is technically feasible, then you will need to conduct the research or studies to determine its feasibility. If still uncertain, then you may need to state what you want as a goal, not as a requirement. Even is a requirement is technically feasible, it may not be attainable due to budget, schedule, or other, e.g., weight, constraints. There is no point in writing a requirement for something you cannot afford -- be reasonable.
- **Clarity** Each requirement should express a single thought, be concise, and simple. It is important that the requirement not be misunderstood -- it must be unambiguous. Simple sentences will most often suffice for a good requirement.

### Common Problems in requirements writing are:

- Making bad assumptions
- Writing implementation (HOW) instead of requirements (WHAT)
- Describing operations instead of writing requirements
- Using incorrect terms
- Using incorrect sentence structure or bad grammar
- Missing requirements
- Over-specifying

### Making bad assumptions

Not enough or not the right amount of information

### Writing implementation (HOW) instead of requirements (WHAT)

There are two major dangers in stating implementation. The one most often cited is that of forcing a design when not intended. The second danger is more subtle and potentially much more detrimental. By stating implementation, the author may be lulled into believing that all requirements are covered. In fact, very important requirements may be missing, and the provider can deliver what as asked for and still not deliver what is wanted.	
**The Implementation Trap**. If you have been doing design studies at a low level, you may begin to document these results as high level requirements -- this is the implementation trap. You will be defining implementation instead of requirements.	

### Describing operations instead of writing requirements

### Using incorrect terms

In a specification, there are terms to be avoided and terms that must be used in a very specific manner. Authors need to understand the use of shall, will, and should:

- **Requirements** use shall.
- **Statements** of fact use will.
- **Goals** use should.

### Using incorrect sentence structure or bad grammar

**Subject Trap.**
**Bad Grammar**. If you use bad grammar you risk that the reader will misinterpret what is stated.

### Unverifiable

Every requirement must be verified.
Ambiguous Term A major cause of unverifiable requirements is the use of ambiguous terms. The terms are ambiguous because they are subjective -- they mean something different to everyone who reads them. This can be avoided by giving people words to avoid. The following lists ambiguous words that we have encountered.

- Minimize
- Maximize
- Rapidq
- User-friendly
- Easy
- Sufficient
- Adequate
- Quick

### Missing requirements

Missing items can be avoided by using a standard outline for your specification, such as those shown in Mil- Std-490 or IEEE P1233, and expanding the outline for your program.
Checklist missing requirements:

- Functional
- Performance
- Interface
- Environment
- Deployment
- Transportation
- Deployment
- Training
- Personnel
- Reliability
- Maintainability
- Operability
- Safety
- Regulatory
- Security
- Privacy
- Design constraints

### Over-specifying

The DoD has stated that over-specification is the primary cause of cost overruns on their programs. Over-specifying is most often from stating something that is unnecessary or from stating overly stringent requirements.
**Unnecessary Items**. Unnecessary requirements creep into a specification in a number of ways. The only cure is careful management review and control.
**Over Stringent**. Most requirements that are too stringent are that way accidentally, not intentionally. A common cause is when an author writes down a number but does not consider the tolerances that are allowable.

[<< Previous Chapter](FreePDM_03-2-SVNProjectStructure.md) | [Content Table](README.md) | [Next Chapter >>](FreePDM_05-Architecture.md)
