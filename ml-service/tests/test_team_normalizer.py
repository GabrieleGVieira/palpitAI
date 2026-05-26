from __future__ import annotations

from app.metrics.schemas import Team
import pandas as pd

from app.metrics.team_normalizer import TeamNormalizer, attach_team_ids


def test_team_normalizer_maps_names_and_aliases() -> None:
    normalizer = TeamNormalizer(
        teams=[Team(id="br", name="Brazil"), Team(id="ci", name="Cote d Ivoire")],
        aliases={"Brasil": "br", "Ivory Coast": "ci"},
    )

    assert normalizer.team_id_for("Brasil") == "br"
    assert normalizer.team_id_for("ivory coast") == "ci"


def test_team_normalizer_reports_unmapped_names() -> None:
    normalizer = TeamNormalizer(teams=[Team(id="br", name="Brazil")], aliases={})

    assert normalizer.team_id_for("Atlantis") is None
    assert normalizer.report_unmapped() == ["Atlantis"]


def test_attach_team_ids_can_keep_unknown_opponents_with_local_ids() -> None:
    normalizer = TeamNormalizer(teams=[Team(id="br", name="Brazil")], aliases={})
    matches = pd.DataFrame([{"home_team": "Brazil", "away_team": "Atlantis"}])

    mapped = attach_team_ids(matches, normalizer, fallback_unknown=True)

    assert mapped.iloc[0]["home_team_id"] == "br"
    assert mapped.iloc[0]["away_team_id"] == "csv:atlantis"
    assert normalizer.report_unmapped() == []
