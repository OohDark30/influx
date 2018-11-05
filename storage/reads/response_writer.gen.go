// Generated by tmpl
// https://github.com/benbjohnson/tmpl
//
// DO NOT EDIT!
// Source: response_writer.gen.go.tmpl

package reads

import (
	"github.com/EMCECS/influx/storage/reads/datatypes"
	"github.com/EMCECS/influx/tsdb/cursors"
)

func (w *ResponseWriter) getFloatPointsFrame() *datatypes.ReadResponse_Frame_FloatPoints {
	var res *datatypes.ReadResponse_Frame_FloatPoints
	if len(w.buffer.Float) > 0 {
		i := len(w.buffer.Float) - 1
		res = w.buffer.Float[i]
		w.buffer.Float[i] = nil
		w.buffer.Float = w.buffer.Float[:i]
	} else {
		res = &datatypes.ReadResponse_Frame_FloatPoints{
			FloatPoints: &datatypes.ReadResponse_FloatPointsFrame{
				Timestamps: make([]int64, 0, batchSize),
				Values:     make([]float64, 0, batchSize),
			},
		}
	}

	return res
}

func (w *ResponseWriter) putFloatPointsFrame(f *datatypes.ReadResponse_Frame_FloatPoints) {
	f.FloatPoints.Timestamps = f.FloatPoints.Timestamps[:0]
	f.FloatPoints.Values = f.FloatPoints.Values[:0]
	w.buffer.Float = append(w.buffer.Float, f)
}

func (w *ResponseWriter) streamFloatArraySeries(cur cursors.FloatArrayCursor) {
	w.sf.DataType = datatypes.DataTypeFloat
	ss := len(w.res.Frames) - 1
	a := cur.Next()
	if len(a.Timestamps) == 0 {
		w.sz -= w.sf.Size()
		w.putSeriesFrame(w.res.Frames[ss].Data.(*datatypes.ReadResponse_Frame_Series))
		w.res.Frames = w.res.Frames[:ss]
	} else if w.sz > writeSize {
		w.Flush()
	}
}

func (w *ResponseWriter) streamFloatArrayPoints(cur cursors.FloatArrayCursor) {
	w.sf.DataType = datatypes.DataTypeFloat
	ss := len(w.res.Frames) - 1

	p := w.getFloatPointsFrame()
	frame := p.FloatPoints
	w.res.Frames = append(w.res.Frames, datatypes.ReadResponse_Frame{Data: p})

	var (
		seriesValueCount = 0
		b                = 0
	)

	for {
		a := cur.Next()
		if len(a.Timestamps) == 0 {
			break
		}

		frame.Timestamps = append(frame.Timestamps, a.Timestamps...)
		frame.Values = append(frame.Values, a.Values...)

		b = len(frame.Timestamps)
		if b >= batchSize {
			seriesValueCount += b
			b = 0
			w.sz += frame.Size()
			if w.sz >= writeSize {
				w.Flush()
				if w.err != nil {
					break
				}
			}

			p = w.getFloatPointsFrame()
			frame = p.FloatPoints
			w.res.Frames = append(w.res.Frames, datatypes.ReadResponse_Frame{Data: p})
		}
	}

	seriesValueCount += b
	w.vc += seriesValueCount
	if seriesValueCount == 0 {
		w.sz -= w.sf.Size()
		w.putSeriesFrame(w.res.Frames[ss].Data.(*datatypes.ReadResponse_Frame_Series))
		w.res.Frames = w.res.Frames[:ss]
	} else if w.sz > writeSize {
		w.Flush()
	}
}

func (w *ResponseWriter) getIntegerPointsFrame() *datatypes.ReadResponse_Frame_IntegerPoints {
	var res *datatypes.ReadResponse_Frame_IntegerPoints
	if len(w.buffer.Integer) > 0 {
		i := len(w.buffer.Integer) - 1
		res = w.buffer.Integer[i]
		w.buffer.Integer[i] = nil
		w.buffer.Integer = w.buffer.Integer[:i]
	} else {
		res = &datatypes.ReadResponse_Frame_IntegerPoints{
			IntegerPoints: &datatypes.ReadResponse_IntegerPointsFrame{
				Timestamps: make([]int64, 0, batchSize),
				Values:     make([]int64, 0, batchSize),
			},
		}
	}

	return res
}

func (w *ResponseWriter) putIntegerPointsFrame(f *datatypes.ReadResponse_Frame_IntegerPoints) {
	f.IntegerPoints.Timestamps = f.IntegerPoints.Timestamps[:0]
	f.IntegerPoints.Values = f.IntegerPoints.Values[:0]
	w.buffer.Integer = append(w.buffer.Integer, f)
}

func (w *ResponseWriter) streamIntegerArraySeries(cur cursors.IntegerArrayCursor) {
	w.sf.DataType = datatypes.DataTypeInteger
	ss := len(w.res.Frames) - 1
	a := cur.Next()
	if len(a.Timestamps) == 0 {
		w.sz -= w.sf.Size()
		w.putSeriesFrame(w.res.Frames[ss].Data.(*datatypes.ReadResponse_Frame_Series))
		w.res.Frames = w.res.Frames[:ss]
	} else if w.sz > writeSize {
		w.Flush()
	}
}

func (w *ResponseWriter) streamIntegerArrayPoints(cur cursors.IntegerArrayCursor) {
	w.sf.DataType = datatypes.DataTypeInteger
	ss := len(w.res.Frames) - 1

	p := w.getIntegerPointsFrame()
	frame := p.IntegerPoints
	w.res.Frames = append(w.res.Frames, datatypes.ReadResponse_Frame{Data: p})

	var (
		seriesValueCount = 0
		b                = 0
	)

	for {
		a := cur.Next()
		if len(a.Timestamps) == 0 {
			break
		}

		frame.Timestamps = append(frame.Timestamps, a.Timestamps...)
		frame.Values = append(frame.Values, a.Values...)

		b = len(frame.Timestamps)
		if b >= batchSize {
			seriesValueCount += b
			b = 0
			w.sz += frame.Size()
			if w.sz >= writeSize {
				w.Flush()
				if w.err != nil {
					break
				}
			}

			p = w.getIntegerPointsFrame()
			frame = p.IntegerPoints
			w.res.Frames = append(w.res.Frames, datatypes.ReadResponse_Frame{Data: p})
		}
	}

	seriesValueCount += b
	w.vc += seriesValueCount
	if seriesValueCount == 0 {
		w.sz -= w.sf.Size()
		w.putSeriesFrame(w.res.Frames[ss].Data.(*datatypes.ReadResponse_Frame_Series))
		w.res.Frames = w.res.Frames[:ss]
	} else if w.sz > writeSize {
		w.Flush()
	}
}

func (w *ResponseWriter) getUnsignedPointsFrame() *datatypes.ReadResponse_Frame_UnsignedPoints {
	var res *datatypes.ReadResponse_Frame_UnsignedPoints
	if len(w.buffer.Unsigned) > 0 {
		i := len(w.buffer.Unsigned) - 1
		res = w.buffer.Unsigned[i]
		w.buffer.Unsigned[i] = nil
		w.buffer.Unsigned = w.buffer.Unsigned[:i]
	} else {
		res = &datatypes.ReadResponse_Frame_UnsignedPoints{
			UnsignedPoints: &datatypes.ReadResponse_UnsignedPointsFrame{
				Timestamps: make([]int64, 0, batchSize),
				Values:     make([]uint64, 0, batchSize),
			},
		}
	}

	return res
}

func (w *ResponseWriter) putUnsignedPointsFrame(f *datatypes.ReadResponse_Frame_UnsignedPoints) {
	f.UnsignedPoints.Timestamps = f.UnsignedPoints.Timestamps[:0]
	f.UnsignedPoints.Values = f.UnsignedPoints.Values[:0]
	w.buffer.Unsigned = append(w.buffer.Unsigned, f)
}

func (w *ResponseWriter) streamUnsignedArraySeries(cur cursors.UnsignedArrayCursor) {
	w.sf.DataType = datatypes.DataTypeUnsigned
	ss := len(w.res.Frames) - 1
	a := cur.Next()
	if len(a.Timestamps) == 0 {
		w.sz -= w.sf.Size()
		w.putSeriesFrame(w.res.Frames[ss].Data.(*datatypes.ReadResponse_Frame_Series))
		w.res.Frames = w.res.Frames[:ss]
	} else if w.sz > writeSize {
		w.Flush()
	}
}

func (w *ResponseWriter) streamUnsignedArrayPoints(cur cursors.UnsignedArrayCursor) {
	w.sf.DataType = datatypes.DataTypeUnsigned
	ss := len(w.res.Frames) - 1

	p := w.getUnsignedPointsFrame()
	frame := p.UnsignedPoints
	w.res.Frames = append(w.res.Frames, datatypes.ReadResponse_Frame{Data: p})

	var (
		seriesValueCount = 0
		b                = 0
	)

	for {
		a := cur.Next()
		if len(a.Timestamps) == 0 {
			break
		}

		frame.Timestamps = append(frame.Timestamps, a.Timestamps...)
		frame.Values = append(frame.Values, a.Values...)

		b = len(frame.Timestamps)
		if b >= batchSize {
			seriesValueCount += b
			b = 0
			w.sz += frame.Size()
			if w.sz >= writeSize {
				w.Flush()
				if w.err != nil {
					break
				}
			}

			p = w.getUnsignedPointsFrame()
			frame = p.UnsignedPoints
			w.res.Frames = append(w.res.Frames, datatypes.ReadResponse_Frame{Data: p})
		}
	}

	seriesValueCount += b
	w.vc += seriesValueCount
	if seriesValueCount == 0 {
		w.sz -= w.sf.Size()
		w.putSeriesFrame(w.res.Frames[ss].Data.(*datatypes.ReadResponse_Frame_Series))
		w.res.Frames = w.res.Frames[:ss]
	} else if w.sz > writeSize {
		w.Flush()
	}
}

func (w *ResponseWriter) getStringPointsFrame() *datatypes.ReadResponse_Frame_StringPoints {
	var res *datatypes.ReadResponse_Frame_StringPoints
	if len(w.buffer.String) > 0 {
		i := len(w.buffer.String) - 1
		res = w.buffer.String[i]
		w.buffer.String[i] = nil
		w.buffer.String = w.buffer.String[:i]
	} else {
		res = &datatypes.ReadResponse_Frame_StringPoints{
			StringPoints: &datatypes.ReadResponse_StringPointsFrame{
				Timestamps: make([]int64, 0, batchSize),
				Values:     make([]string, 0, batchSize),
			},
		}
	}

	return res
}

func (w *ResponseWriter) putStringPointsFrame(f *datatypes.ReadResponse_Frame_StringPoints) {
	f.StringPoints.Timestamps = f.StringPoints.Timestamps[:0]
	f.StringPoints.Values = f.StringPoints.Values[:0]
	w.buffer.String = append(w.buffer.String, f)
}

func (w *ResponseWriter) streamStringArraySeries(cur cursors.StringArrayCursor) {
	w.sf.DataType = datatypes.DataTypeString
	ss := len(w.res.Frames) - 1
	a := cur.Next()
	if len(a.Timestamps) == 0 {
		w.sz -= w.sf.Size()
		w.putSeriesFrame(w.res.Frames[ss].Data.(*datatypes.ReadResponse_Frame_Series))
		w.res.Frames = w.res.Frames[:ss]
	} else if w.sz > writeSize {
		w.Flush()
	}
}

func (w *ResponseWriter) streamStringArrayPoints(cur cursors.StringArrayCursor) {
	w.sf.DataType = datatypes.DataTypeString
	ss := len(w.res.Frames) - 1

	p := w.getStringPointsFrame()
	frame := p.StringPoints
	w.res.Frames = append(w.res.Frames, datatypes.ReadResponse_Frame{Data: p})

	var (
		seriesValueCount = 0
		b                = 0
	)

	for {
		a := cur.Next()
		if len(a.Timestamps) == 0 {
			break
		}

		frame.Timestamps = append(frame.Timestamps, a.Timestamps...)
		frame.Values = append(frame.Values, a.Values...)

		b = len(frame.Timestamps)
		if b >= batchSize {
			seriesValueCount += b
			b = 0
			w.sz += frame.Size()
			if w.sz >= writeSize {
				w.Flush()
				if w.err != nil {
					break
				}
			}

			p = w.getStringPointsFrame()
			frame = p.StringPoints
			w.res.Frames = append(w.res.Frames, datatypes.ReadResponse_Frame{Data: p})
		}
	}

	seriesValueCount += b
	w.vc += seriesValueCount
	if seriesValueCount == 0 {
		w.sz -= w.sf.Size()
		w.putSeriesFrame(w.res.Frames[ss].Data.(*datatypes.ReadResponse_Frame_Series))
		w.res.Frames = w.res.Frames[:ss]
	} else if w.sz > writeSize {
		w.Flush()
	}
}

func (w *ResponseWriter) getBooleanPointsFrame() *datatypes.ReadResponse_Frame_BooleanPoints {
	var res *datatypes.ReadResponse_Frame_BooleanPoints
	if len(w.buffer.Boolean) > 0 {
		i := len(w.buffer.Boolean) - 1
		res = w.buffer.Boolean[i]
		w.buffer.Boolean[i] = nil
		w.buffer.Boolean = w.buffer.Boolean[:i]
	} else {
		res = &datatypes.ReadResponse_Frame_BooleanPoints{
			BooleanPoints: &datatypes.ReadResponse_BooleanPointsFrame{
				Timestamps: make([]int64, 0, batchSize),
				Values:     make([]bool, 0, batchSize),
			},
		}
	}

	return res
}

func (w *ResponseWriter) putBooleanPointsFrame(f *datatypes.ReadResponse_Frame_BooleanPoints) {
	f.BooleanPoints.Timestamps = f.BooleanPoints.Timestamps[:0]
	f.BooleanPoints.Values = f.BooleanPoints.Values[:0]
	w.buffer.Boolean = append(w.buffer.Boolean, f)
}

func (w *ResponseWriter) streamBooleanArraySeries(cur cursors.BooleanArrayCursor) {
	w.sf.DataType = datatypes.DataTypeBoolean
	ss := len(w.res.Frames) - 1
	a := cur.Next()
	if len(a.Timestamps) == 0 {
		w.sz -= w.sf.Size()
		w.putSeriesFrame(w.res.Frames[ss].Data.(*datatypes.ReadResponse_Frame_Series))
		w.res.Frames = w.res.Frames[:ss]
	} else if w.sz > writeSize {
		w.Flush()
	}
}

func (w *ResponseWriter) streamBooleanArrayPoints(cur cursors.BooleanArrayCursor) {
	w.sf.DataType = datatypes.DataTypeBoolean
	ss := len(w.res.Frames) - 1

	p := w.getBooleanPointsFrame()
	frame := p.BooleanPoints
	w.res.Frames = append(w.res.Frames, datatypes.ReadResponse_Frame{Data: p})

	var (
		seriesValueCount = 0
		b                = 0
	)

	for {
		a := cur.Next()
		if len(a.Timestamps) == 0 {
			break
		}

		frame.Timestamps = append(frame.Timestamps, a.Timestamps...)
		frame.Values = append(frame.Values, a.Values...)

		b = len(frame.Timestamps)
		if b >= batchSize {
			seriesValueCount += b
			b = 0
			w.sz += frame.Size()
			if w.sz >= writeSize {
				w.Flush()
				if w.err != nil {
					break
				}
			}

			p = w.getBooleanPointsFrame()
			frame = p.BooleanPoints
			w.res.Frames = append(w.res.Frames, datatypes.ReadResponse_Frame{Data: p})
		}
	}

	seriesValueCount += b
	w.vc += seriesValueCount
	if seriesValueCount == 0 {
		w.sz -= w.sf.Size()
		w.putSeriesFrame(w.res.Frames[ss].Data.(*datatypes.ReadResponse_Frame_Series))
		w.res.Frames = w.res.Frames[:ss]
	} else if w.sz > writeSize {
		w.Flush()
	}
}