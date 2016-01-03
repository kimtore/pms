.md.1:
	@PANDOC@ -s -t man $< -o $@

man1_MANS = pms.1

mostlyclean-local:
	-rm $(man1_MANS)

EXTRA_DIST = $(man1_MANS)
