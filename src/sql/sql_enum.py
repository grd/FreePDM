#!/usr/bin python3
# -*- coding: utf-8 -*-
"""
    :copyright: Copyright 2022 by the FreePDM team
    :license:   MIT License.
"""

import enum


class ProjectState(enum.Enum):
    """Class with enumeration of the Project state"""

    # https://www.indeed.com/career-advice/career-development/project-statuses
    upcoming = "Upcoming"
    pending = "Pending"
    notstarted = "Not Started"
    draft = "Draft"
    active = "Active"
    priority = "Priority"
    canceld = "Canceled"
    onhold = "On-Hold"
    archived = "Archived"
    # check also: https://support.ptc.com/help/wnc/r11.2.0.0/en/index.html#page/Windchill_Help_Center/ProjMgmtPhaseState.html


class DensityUnit(enum.Enum):
    """Class with enumeration of mass units"""

    gmm3 = "Gram / mm^3"
    gcm3 = "Gram / cm^3"
    kgm3 = "Kilo Gram / m^3"
    tonnem3 = "Tonne / m^3"
    metrictonm3 = "Metric Ton / m^3"  # equal to tonne
    poundft3 = "Pound / Ft^3"
    poundinch3 = "Pound / inch^3"


class VolumeUnit(enum.Enum):
    """Class with enumeration of mass units"""

    mm3 = "Cubic mm [mm^3]"
    cm3 = "Cubic cm [cm^3]"
    m3 = "Cubic m [m^3]"
    litre = "Litre"
    ft3 = "Cubic ft [Ft^3]"
    inch3 = "Cubic inch [Inch^3]"


class WeightUnit(enum.Enum):
    """Class with enumeration of mass units"""

    g = "Gram [g]"
    kg = "Kilo Gram [kg]"
    tonne = "Tonne [t]"
    metricton = "Metric Ton"  # equal to tonne
    pound = "Pound [p]"
    slug = "Slug"


class AreaUnit(enum.Enum):
    """Class with enumeration of mass units"""

    mm2 = "Square mm [mm^2]"
    cm2 = "Square cm [cm^2]"
    m2 = "Square m [m^2]"
    ft2 = "Square ft [Ft^2]"
    inch2 = "Square inch [Inch^2]"


class RevisionState(enum.Enum):
    """Class with enumeration of the Project state"""

    concept = "Concept"
    underreview = "Under Review"
    released = "Released"
    inwork = "In-Work"
    depreciated = "Depreciated"


class TracebilityState(enum.Enum):
    """Class with enumeration of the Project state"""

    lot = "Lot"
    lotserial = "Lot and Serial"
    serieal = "Serial"
    nottraced = "Not-Traced"
